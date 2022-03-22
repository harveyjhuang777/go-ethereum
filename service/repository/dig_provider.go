package repository

import (
	"sync"

	"go.uber.org/dig"
	"gorm.io/gorm"

	"github.com/harveyjhuang777/go-ethereum/service/model"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/dbcli"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

var (
	once sync.Once
	self *packet
)

func NewRepository(in digIn) digOut {
	once.Do(func() {
		self = &packet{
			in: in,
			digOut: digOut{
				BlockRepository:       newBlockRepository(in),
				TransactionRepository: newTransactionRepository(in),
			},
		}
	})

	return self.digOut
}

func Migration(db *gorm.DB) error {
	// Migration
	if err := db.AutoMigrate(&model.Block{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.Transaction{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.TransactionLog{}); err != nil {
		return err
	}

	return nil
}

type digIn struct {
	dig.In

	DB     dbcli.IMySQLClient
	Logger logger.ILogger
}

type packet struct {
	in digIn

	digOut
}

type digOut struct {
	dig.Out

	BlockRepository       IBlock
	TransactionRepository ITransaction
}

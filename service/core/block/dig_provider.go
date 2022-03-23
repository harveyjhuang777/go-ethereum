package block

import (
	"sync"

	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/repository"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/dbcli"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/snowflake"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

var (
	once sync.Once
	self *packet
)

func NewBlock(in digIn) digOut {
	once.Do(func() {
		self = &packet{
			in: in,
			digOut: digOut{
				BlockListUseCase: newBlockList(in),
			},
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	DB          dbcli.IMySQLClient
	IdGenerator snowflake.IIDGenerator
	Logger      logger.ILogger

	BlockRepository repository.IBlock
}

type packet struct {
	in digIn

	digOut
}

type digOut struct {
	dig.Out

	BlockListUseCase IBlockList
}

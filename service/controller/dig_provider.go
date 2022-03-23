package controller

import (
	"sync"

	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/core/block"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

var (
	once sync.Once
	self *packet
)

func NewRestController(in digIn) digOut {
	once.Do(func() {
		self = &packet{
			in: in,
			digOut: digOut{
				BlockController: newBlockController(in),
			},
		}

	})

	return self.digOut
}

type packet struct {
	in digIn

	digOut
}

type digIn struct {
	dig.In

	Logger                   logger.ILogger
	BlockListUseCase         block.IBlockList
	BlockDetailUseCase       block.IBlockDetail
	TransactionDetailUseCase block.ITransactionDetail
}

type digOut struct {
	dig.Out

	BlockController IBlockController
}

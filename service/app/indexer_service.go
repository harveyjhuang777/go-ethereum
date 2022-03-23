package app

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/core/block"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/dbcli"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/ethcli"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/snowflake"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

var (
	indexOnce sync.Once
	indexSelf *indexerPacket
)

func NewIndexerService(in indexerIn) indexerOut {
	indexOnce.Do(func() {
		indexSelf = &indexerPacket{
			in: in,
			indexerOut: indexerOut{
				IndexerService: newIndexerService(in),
			},
		}
	})

	return indexSelf.indexerOut
}

type indexerIn struct {
	dig.In

	DB          dbcli.IMySQLClient
	IdGenerator snowflake.IIDGenerator
	Logger      logger.ILogger
	EthApiCli   ethcli.IEthCli

	BlockInsertUseCase block.IBlockInsert
}

type indexerPacket struct {
	in indexerIn

	indexerOut
}

type indexerOut struct {
	dig.Out

	IndexerService IIndexerService
}

type IIndexerService interface {
	Run(ctx context.Context)
}

func newIndexerService(in indexerIn) IIndexerService {
	return &indexerService{in: in}
}

type indexerService struct {
	in indexerIn
}

func (srv *indexerService) Run(ctx context.Context) {
	headers := make(chan *types.Header)
	sub, err := srv.in.EthApiCli.SubscribeNewHead(ctx, headers)
	if err != nil {
		srv.in.Logger.Panic(ctx, err)
	}
	srv.in.BlockInsertUseCase.Handle(ctx, sub, headers)
}

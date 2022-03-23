package ethcli

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/repository"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/dbcli"
	"github.com/harveyjhuang777/go-ethereum/service/util/config"
	"github.com/harveyjhuang777/go-ethereum/service/util/filepath"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

type ethApiTestApp struct {
	dig.In

	EthApi IEthCli
	DB     dbcli.IMySQLClient
}

type ethApiTestSuite struct {
	suite.Suite
	ctx context.Context
	app *ethApiTestApp
}

func (s *ethApiTestSuite) SetupSuite() {
	filepath.InitRootFolder("../../..")
	// ctx
	s.ctx = context.Background()

	binder := dig.New()
	s.Require().Nil(binder.Provide(config.NewConfig))
	s.Require().Nil(binder.Provide(logger.NewSysLog))
	s.Require().Nil(binder.Provide(dbcli.NewDBClient))
	s.Require().Nil(binder.Provide(NewEthCli))
	s.Require().Nil(binder.Invoke(func(app ethApiTestApp) {
		s.app = &app
	}))

	// clear table
	s.Require().Nil(repository.Migration(s.app.DB.Session()))
}

func (s *ethApiTestSuite) SetupTest() {

}

func (s *ethApiTestSuite) TearDownTest() {
}

func (s *ethApiTestSuite) TestEthCli() {
	header, err := s.app.EthApi.GetLatestHeader(s.ctx)
	s.Require().Nil(err)
	s.Require().NotEmpty(header)

	block, err := s.app.EthApi.GetBlockByNumber(s.ctx, big.NewInt(17536859))
	s.Require().Nil(err)
	s.Require().NotEmpty(block)

	txn, _, err := s.app.EthApi.GetTransactionByHash(s.ctx, block.Transactions()[0].Hash())
	s.Require().Nil(err)
	s.Require().NotEmpty(txn)

	receipt, err := s.app.EthApi.GetTransactionReceipt(s.ctx, txn.Hash())
	s.Require().Nil(err)
	s.Require().NotEmpty(receipt)

	client, err := ethclient.Dial("wss://bsc-ws-node.nariox.org:443")
	s.Require().Nil(err)

	s.app.EthApi = &ethCli{client: client}
	headers := make(chan *types.Header)
	sub, err := s.app.EthApi.SubscribeNewHead(s.ctx, headers)
	s.Require().Nil(err)
	s.Require().NotEmpty(sub)
}

func TestEthApiRepository(t *testing.T) {
	suite.Run(t, &ethApiTestSuite{})
}

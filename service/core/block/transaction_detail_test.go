package block

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"

	"github.com/harveyjhuang777/go-ethereum/service/model"
	"github.com/harveyjhuang777/go-ethereum/service/repository"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/dbcli"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/ethcli"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/snowflake"
	"github.com/harveyjhuang777/go-ethereum/service/util/codebook"
	"github.com/harveyjhuang777/go-ethereum/service/util/config"
	"github.com/harveyjhuang777/go-ethereum/service/util/filepath"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

type transactionDetailTestApp struct {
	dig.In

	UseCase ITransactionDetail
	DB      dbcli.IMySQLClient
}

type transactionDetailSuite struct {
	suite.Suite
	ctx context.Context
	app *transactionDetailTestApp
}

func (s *transactionDetailSuite) SetupSuite() {
	filepath.InitRootFolder("../../..")
	// ctx
	s.ctx = context.Background()

	binder := dig.New()
	s.Require().Nil(binder.Provide(newTransactionDetail))
	s.Require().Nil(binder.Provide(config.NewConfig))
	s.Require().Nil(binder.Provide(dbcli.NewDBClient))
	s.Require().Nil(binder.Provide(logger.NewSysLog))
	s.Require().Nil(binder.Provide(ethcli.NewEthCli))
	s.Require().Nil(binder.Provide(snowflake.NewIDGenerator))
	s.Require().Nil(binder.Provide(repository.NewRepository))
	s.Require().Nil(binder.Invoke(func(app transactionDetailTestApp) {
		s.app = &app
	}))

	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.TransactionLog{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Transaction{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *transactionDetailSuite) SetupTest() {
	now := time.Now().UTC()
	basicTestCase1 := &model.Block{
		Number:     1,
		Hash:       "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae",
		Time:       uint64(now.Unix()),
		ParentHash: "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
	}

	// Test insert
	s.Require().Nil(s.app.DB.Session().Create(basicTestCase1).Error)

	basicTestCase2 := &model.Transaction{
		BlockNumber: basicTestCase1.Number,
		Hash:        "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
		From:        "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
		To:          "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
		Nonce:       21,
		Value:       4290000000000000,
	}

	// Test insert
	s.Require().Nil(s.app.DB.Session().Create(basicTestCase2).Error)

	basicTestCase := &model.TransactionLog{
		ID:              1,
		TransactionHash: basicTestCase2.Hash,
		Index:           1,
		Data:            "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
	}

	// Test insert
	s.Require().Nil(s.app.DB.Session().Create(basicTestCase).Error)
}

func (s *transactionDetailSuite) TearDownTest() {
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.TransactionLog{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Transaction{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *transactionDetailSuite) TestDetail() {
	testCase1 := &model.Transaction{
		BlockNumber: 1,
		Hash:        "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df1",
		From:        "0xa7d9ddbe1f17865597fbd27ec712455208b6b16d",
		To:          "0xf02c1c8e6114b1dbe8937a39260b5b0a374433bb",
		Nonce:       22,
		Value:       4190000000000000,
	}

	s.Require().Empty(s.app.DB.Session().Create(testCase1).Error)

	resp, err := s.app.UseCase.Handle(s.ctx, "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2")
	s.Require().Empty(err)
	s.Require().Equal(1, len(resp.Logs))
	s.Require().Equal("0xa7d9ddbe1f17865597fbd27ec712455208b6b76d", resp.From)
	s.Require().Equal("0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb", resp.To)
	s.Require().EqualValues(21, resp.Nonce)
	s.Require().EqualValues(4290000000000000, resp.Value)

	resp, err = s.app.UseCase.Handle(s.ctx, testCase1.Hash)
	s.Require().Empty(err)
	s.Require().Equal(0, len(resp.Logs))
	s.Require().Equal(testCase1.From, resp.From)
	s.Require().Equal(testCase1.To, resp.To)
	s.Require().Equal(testCase1.Nonce, resp.Nonce)
	s.Require().Equal(testCase1.Value, resp.Value)

	resp, err = s.app.UseCase.Handle(s.ctx, "")
	s.Require().ErrorIs(codebook.ErrDataNotExist, err)
}

func TestTransactionDetail(t *testing.T) {
	suite.Run(t, &transactionDetailSuite{})
}

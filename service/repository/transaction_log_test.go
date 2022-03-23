package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"
	"gorm.io/gorm"

	"github.com/harveyjhuang777/go-ethereum/service/model"
	"github.com/harveyjhuang777/go-ethereum/service/thirdparty/dbcli"
	"github.com/harveyjhuang777/go-ethereum/service/util/config"
	"github.com/harveyjhuang777/go-ethereum/service/util/filepath"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

type transactionLogTestApp struct {
	dig.In

	TransactionLogRepository ITransactionLog
	DB                       dbcli.IMySQLClient
}

type transactionLogTestSuite struct {
	suite.Suite
	ctx context.Context
	app *transactionLogTestApp
}

func (s *transactionLogTestSuite) SetupSuite() {
	filepath.InitRootFolder("../../..")
	// ctx
	s.ctx = context.Background()

	binder := dig.New()
	s.Require().Nil(binder.Provide(config.NewConfig))
	s.Require().Nil(binder.Provide(dbcli.NewDBClient))
	s.Require().Nil(binder.Provide(logger.NewSysLog))
	s.Require().Nil(binder.Provide(NewRepository))
	s.Require().Nil(binder.Invoke(func(app transactionLogTestApp) {
		s.app = &app
	}))

	// clear table
	s.Require().Nil(Migration(s.app.DB.Session()))
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.TransactionLog{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Transaction{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *transactionLogTestSuite) SetupTest() {
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

func (s *transactionLogTestSuite) TearDownTest() {
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.TransactionLog{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Transaction{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *transactionLogTestSuite) TestInsert() {
	// Test duplicate insert on PK
	testCase1 := &model.TransactionLog{
		ID:              1,
		TransactionHash: "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
		Index:           1,
		Data:            "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
	}

	err := s.app.TransactionLogRepository.Insert(s.ctx, s.app.DB.Session(), testCase1)
	s.Require().NotEmpty(err)
}

func (s *transactionLogTestSuite) TestList() {
	testCase1 := &model.TransactionLog{
		ID:              2,
		TransactionHash: "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
		Index:           2,
		Data:            "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df3",
	}

	s.Require().Empty(s.app.TransactionLogRepository.Insert(s.ctx, s.app.DB.Session(), testCase1))

	resp, err := s.app.TransactionLogRepository.List(s.ctx, s.app.DB.Session())
	s.Require().Empty(err)
	s.Require().Equal(2, len(resp))

	condFunc := func(db *gorm.DB) *gorm.DB {
		return db.Limit(1).Order("created_at desc")
	}
	resp, err = s.app.TransactionLogRepository.List(s.ctx, s.app.DB.Session(), condFunc)
	s.Require().Empty(err)
	s.Require().Equal(1, len(resp))
	s.Require().EqualValues(2, resp[0].ID)
}

func TestTransactionLogRepository(t *testing.T) {
	suite.Run(t, &transactionLogTestSuite{})
}

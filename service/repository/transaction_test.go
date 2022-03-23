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

type transactionTestApp struct {
	dig.In

	TransactionRepository ITransaction
	DB                    dbcli.IMySQLClient
}

type transactionTestSuite struct {
	suite.Suite
	ctx context.Context
	app *transactionTestApp
}

func (s *transactionTestSuite) SetupSuite() {
	filepath.InitRootFolder("../../..")
	// ctx
	s.ctx = context.Background()

	binder := dig.New()
	s.Require().Nil(binder.Provide(config.NewConfig))
	s.Require().Nil(binder.Provide(dbcli.NewDBClient))
	s.Require().Nil(binder.Provide(logger.NewSysLog))
	s.Require().Nil(binder.Provide(NewRepository))
	s.Require().Nil(binder.Invoke(func(app transactionTestApp) {
		s.app = &app
	}))

	// clear table
	s.Require().Nil(Migration(s.app.DB.Session()))
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Transaction{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *transactionTestSuite) SetupTest() {
	now := time.Now().UTC()
	number := 436
	hash := "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae"
	basicTestCase1 := &model.Block{
		ID:         1,
		Time:       now.Unix(),
		ParentHash: "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
	}
	basicTestCase1.Number = &number
	basicTestCase1.Hash = &hash

	// Test insert
	s.Require().Nil(s.app.DB.Session().Create(basicTestCase1).Error)

	basicTestCase := &model.Transaction{
		ID:      1,
		BlockID: 1,
		Hash:    "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
		From:    "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
		To:      "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
		Nonce:   21,
		Value:   4290000000000000,
	}

	// Test insert
	s.Require().Nil(s.app.DB.Session().Create(basicTestCase).Error)
}

func (s *transactionTestSuite) TearDownTest() {
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Transaction{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *transactionTestSuite) TestInsert() {
	// Test duplicate insert on PK
	testCase1 := &model.Transaction{
		ID:      1,
		BlockID: 1,
		Hash:    "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
		From:    "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
		To:      "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
		Nonce:   21,
		Value:   4290000000000000,
	}

	err := s.app.TransactionRepository.Insert(s.ctx, s.app.DB.Session(), testCase1)
	s.Require().NotEmpty(err)
}

func (s *transactionTestSuite) TestUpdate() {

	testCase1, err := s.app.TransactionRepository.First(s.ctx, s.app.DB.Session(), 1)
	s.Require().Nil(err)

	testCase1.Hash = "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21bc"
	s.Require().Nil(s.app.TransactionRepository.Update(s.ctx, s.app.DB.Session(), testCase1))

	resp, err := s.app.TransactionRepository.First(s.ctx, s.app.DB.Session(), 1)
	s.Require().Nil(err)
	s.Require().Equal(resp.Hash, testCase1.Hash)
	s.Require().Equal(resp.CreatedAt, testCase1.CreatedAt)
	s.Require().NotEqual(resp.UpdatedAt, testCase1.UpdatedAt)
}

func (s *transactionTestSuite) TestFirst() {
	expected := &model.Transaction{
		ID:      1,
		BlockID: 1,
		Hash:    "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
		From:    "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
		To:      "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
		Nonce:   21,
		Value:   4290000000000000,
	}

	resp, err := s.app.TransactionRepository.First(s.ctx, s.app.DB.Session(), 1)
	// Test first
	s.Require().Nil(err)
	s.Require().EqualValues(expected.From, resp.From)
	s.Require().EqualValues(expected.Hash, resp.Hash)
	s.Require().EqualValues(expected.Value, resp.Value)

	// Test Record Not Found
	resp, err = s.app.TransactionRepository.First(s.ctx, s.app.DB.Session(), 2)
	s.Require().Nil(resp)
	s.Require().ErrorIs(gorm.ErrRecordNotFound, err)
}

func (s *transactionTestSuite) TestList() {
	testCase1 := &model.Transaction{
		ID:      2,
		BlockID: 1,
		Hash:    "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df3",
		From:    "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
		To:      "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
		Nonce:   22,
		Value:   4290000000000000,
	}

	s.Require().Empty(s.app.TransactionRepository.Insert(s.ctx, s.app.DB.Session(), testCase1))

	resp, err := s.app.TransactionRepository.List(s.ctx, s.app.DB.Session())
	s.Require().Empty(err)
	s.Require().Equal(2, len(resp))

	condFunc := func(db *gorm.DB) *gorm.DB {
		return db.Limit(1).Order("created_at desc")
	}
	resp, err = s.app.TransactionRepository.List(s.ctx, s.app.DB.Session(), condFunc)
	s.Require().Empty(err)
	s.Require().Equal(1, len(resp))
	s.Require().EqualValues(2, resp[0].ID)
}

func (s *transactionTestSuite) TestFirstByHash() {
	expected := &model.Transaction{
		ID:      1,
		BlockID: 1,
		Hash:    "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
		From:    "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
		To:      "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
		Nonce:   21,
		Value:   4290000000000000,
	}

	resp, err := s.app.TransactionRepository.FirstByHash(s.ctx, s.app.DB.Session(), expected.Hash)
	// Test first
	s.Require().Nil(err)
	s.Require().EqualValues(expected.From, resp.From)
	s.Require().EqualValues(expected.Hash, resp.Hash)
	s.Require().EqualValues(expected.Value, resp.Value)

	// Test Record Not Found
	resp, err = s.app.TransactionRepository.FirstByHash(s.ctx, s.app.DB.Session(), "")
	s.Require().Nil(resp)
	s.Require().ErrorIs(gorm.ErrRecordNotFound, err)
}

func TestTransactionRepository(t *testing.T) {
	suite.Run(t, &transactionTestSuite{})
}

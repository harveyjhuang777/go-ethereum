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

type blockTestApp struct {
	dig.In

	BlockRepository IBlock
	DB              dbcli.IMySQLClient
}

type blockTestSuite struct {
	suite.Suite
	ctx context.Context
	app *blockTestApp
}

func (s *blockTestSuite) SetupSuite() {
	filepath.InitRootFolder("../../..")
	// ctx
	s.ctx = context.Background()

	binder := dig.New()
	s.Require().Nil(binder.Provide(config.NewConfig))
	s.Require().Nil(binder.Provide(dbcli.NewDBClient))
	s.Require().Nil(binder.Provide(logger.NewSysLog))
	s.Require().Nil(binder.Provide(NewRepository))
	s.Require().Nil(binder.Invoke(func(app blockTestApp) {
		s.app = &app
	}))

	// clear table
	s.Require().Nil(Migration(s.app.DB.Session()))
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *blockTestSuite) SetupTest() {
	now := time.Now().UTC()
	number := 436
	hash := "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae"
	basicTestCase := &model.Block{
		ID:         1,
		Time:       now.Unix(),
		ParentHash: "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
	}
	basicTestCase.Number = &number
	basicTestCase.Hash = &hash

	// Test insert
	s.Require().Nil(s.app.DB.Session().Create(basicTestCase).Error)
}

func (s *blockTestSuite) TearDownTest() {
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *blockTestSuite) TestInsert() {
	// Test duplicate insert on PK
	now := time.Now().UTC()
	number := 436
	hash := "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae"
	testCase1 := &model.Block{
		ID:         1,
		Time:       now.Unix(),
		ParentHash: "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
	}
	testCase1.Number = &number
	testCase1.Hash = &hash

	err := s.app.BlockRepository.Insert(s.ctx, s.app.DB.Session(), testCase1)
	s.Require().NotEmpty(err)
}

func (s *blockTestSuite) TestUpdate() {

	testCase1, err := s.app.BlockRepository.First(s.ctx, s.app.DB.Session(), 1)
	s.Require().Nil(err)

	*testCase1.Hash = "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21bc"
	s.Require().Nil(s.app.BlockRepository.Update(s.ctx, s.app.DB.Session(), testCase1))

	resp, err := s.app.BlockRepository.First(s.ctx, s.app.DB.Session(), 1)
	s.Require().Nil(err)
	s.Require().Equal(resp.Hash, testCase1.Hash)
	s.Require().Equal(resp.CreatedAt, testCase1.CreatedAt)
	s.Require().NotEqual(resp.UpdatedAt, testCase1.UpdatedAt)
}

func (s *blockTestSuite) TestFirst() {
	now := time.Now().UTC()
	number := 436
	hash := "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae"
	expected := &model.Block{
		ID:         1,
		Time:       now.Unix(),
		ParentHash: "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
	}
	expected.Number = &number
	expected.Hash = &hash

	resp, err := s.app.BlockRepository.First(s.ctx, s.app.DB.Session(), 1)
	// Test first
	s.Require().Nil(err)
	s.Require().EqualValues(expected.Number, resp.Number)
	s.Require().EqualValues(expected.Hash, resp.Hash)
	s.Require().EqualValues(expected.ParentHash, resp.ParentHash)

	// Test Record Not Found
	resp, err = s.app.BlockRepository.First(s.ctx, s.app.DB.Session(), 2)
	s.Require().Nil(resp)
	s.Require().ErrorIs(gorm.ErrRecordNotFound, err)
}

func (s *blockTestSuite) TestList() {
	now := time.Now().UTC()
	number := 437
	hash := "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21cd"
	testCase1 := &model.Block{
		ID:         2,
		Time:       now.Unix(),
		ParentHash: "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a53",
	}
	testCase1.Number = &number
	testCase1.Hash = &hash

	s.Require().Empty(s.app.BlockRepository.Insert(s.ctx, s.app.DB.Session(), testCase1))

	resp, err := s.app.BlockRepository.List(s.ctx, s.app.DB.Session())
	s.Require().Empty(err)
	s.Require().Equal(2, len(resp))

	condFunc := func(db *gorm.DB) *gorm.DB {
		return db.Limit(1).Order("created_at desc")
	}
	resp, err = s.app.BlockRepository.List(s.ctx, s.app.DB.Session(), condFunc)
	s.Require().Empty(err)
	s.Require().Equal(1, len(resp))
	s.Require().EqualValues(2, resp[0].ID)
}

func TestBlockRepository(t *testing.T) {
	suite.Run(t, &blockTestSuite{})
}

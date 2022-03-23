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
	"github.com/harveyjhuang777/go-ethereum/service/util/config"
	"github.com/harveyjhuang777/go-ethereum/service/util/filepath"
	"github.com/harveyjhuang777/go-ethereum/service/util/logger"
)

type blockListTestApp struct {
	dig.In

	UseCase IBlockList
	DB      dbcli.IMySQLClient
}

type blockListSuite struct {
	suite.Suite
	ctx context.Context
	app *blockListTestApp
}

func (s *blockListSuite) SetupSuite() {
	filepath.InitRootFolder("../../..")
	// ctx
	s.ctx = context.Background()

	binder := dig.New()
	s.Require().Nil(binder.Provide(newBlockList))
	s.Require().Nil(binder.Provide(config.NewConfig))
	s.Require().Nil(binder.Provide(dbcli.NewDBClient))
	s.Require().Nil(binder.Provide(logger.NewSysLog))
	s.Require().Nil(binder.Provide(ethcli.NewEthCli))
	s.Require().Nil(binder.Provide(snowflake.NewIDGenerator))
	s.Require().Nil(binder.Provide(repository.NewRepository))
	s.Require().Nil(binder.Invoke(func(app blockListTestApp) {
		s.app = &app
	}))

	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.TransactionLog{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Transaction{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *blockListSuite) SetupTest() {
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

func (s *blockListSuite) TearDownTest() {
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.TransactionLog{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Transaction{}).Error)
	s.Require().Nil(s.app.DB.Session().Where("1 = 1").Delete(&model.Block{}).Error)
}

func (s *blockListSuite) TestList() {
	now := time.Now().UTC()
	testCase1 := &model.Block{
		Number:     2,
		Hash:       "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21cd",
		Time:       uint64(now.Unix()),
		ParentHash: "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a53",
	}

	s.Require().Empty(s.app.DB.Session().Create(testCase1).Error)

	resp, err := s.app.UseCase.Handle(s.ctx, 10)
	s.Require().Empty(err)
	s.Require().Equal(2, len(resp.Blocks))
	s.Require().EqualValues(testCase1.Number, resp.Blocks[0].Number)
	s.Require().EqualValues(testCase1.Hash, resp.Blocks[0].Hash)
	s.Require().EqualValues(testCase1.ParentHash, resp.Blocks[0].ParentHash)

	resp, err = s.app.UseCase.Handle(s.ctx, 1)
	s.Require().Empty(err)
	s.Require().Equal(1, len(resp.Blocks))
	s.Require().EqualValues(testCase1.Number, resp.Blocks[0].Number)
	s.Require().EqualValues(testCase1.Hash, resp.Blocks[0].Hash)
	s.Require().EqualValues(testCase1.ParentHash, resp.Blocks[0].ParentHash)

	resp, err = s.app.UseCase.Handle(s.ctx, 0)
	s.Require().Empty(err)
	s.Require().Equal(2, len(resp.Blocks))
	s.Require().EqualValues(testCase1.Number, resp.Blocks[0].Number)
	s.Require().EqualValues(testCase1.Hash, resp.Blocks[0].Hash)
	s.Require().EqualValues(testCase1.ParentHash, resp.Blocks[0].ParentHash)
}

func TestBlockList(t *testing.T) {
	suite.Run(t, &blockListSuite{})
}

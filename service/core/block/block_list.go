package block

import (
	"context"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/harveyjhuang777/go-ethereum/service/model"
	"github.com/harveyjhuang777/go-ethereum/service/util/codebook"
)

type IBlockList interface {
	Handle(ctx context.Context, limit int) (*model.BlockListResp, error)
}

func newBlockList(in digIn) IBlockList {
	return &blockList{
		in: in,
	}
}

type blockList struct {
	in digIn
}

func (uc *blockList) Handle(ctx context.Context, limit int) (*model.BlockListResp, error) {
	db := uc.in.DB.Session()

	condFuncs := uc.generateCondFunc(limit)

	blocks, err := uc.in.BlockRepository.List(ctx, db, condFuncs...)
	if err != nil {
		uc.in.Logger.Error(ctx, err)
		return nil, codebook.ErrDatabase
	}

	var blockList []*model.BlockList

	if err := copier.Copy(&blockList, &blocks); err != nil {
		uc.in.Logger.Error(ctx, err)
		return nil, codebook.ErrServer
	}

	return &model.BlockListResp{Blocks: blockList}, nil
}

func (uc *blockList) generateCondFunc(limit int) []func(*gorm.DB) *gorm.DB {
	var condFuncs []func(*gorm.DB) *gorm.DB

	sortFunc := func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}

	condFuncs = append(condFuncs, sortFunc)

	if limit > 0 {
		condFunc := func(db *gorm.DB) *gorm.DB {
			return db.Limit(limit)
		}
		condFuncs = append(condFuncs, condFunc)
	}
	return condFuncs
}

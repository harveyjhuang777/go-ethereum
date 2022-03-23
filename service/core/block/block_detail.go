package block

import (
	"context"

	"github.com/jinzhu/copier"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/harveyjhuang777/go-ethereum/service/model"
	"github.com/harveyjhuang777/go-ethereum/service/util/codebook"
)

type IBlockDetail interface {
	Handle(ctx context.Context, id int64) (*model.BlockDetail, error)
}

func newBlockDetail(in digIn) IBlockDetail {
	return &blockDetail{
		in: in,
	}
}

type blockDetail struct {
	in digIn
}

func (uc *blockDetail) Handle(ctx context.Context, id int64) (*model.BlockDetail, error) {
	db := uc.in.DB.Session()

	block, err := uc.in.BlockRepository.First(ctx, db, id)
	if err != nil {
		if xerrors.Is(gorm.ErrRecordNotFound, err) {
			return nil, codebook.ErrDataNotExist
		}
		uc.in.Logger.Error(ctx, err)
		return nil, codebook.ErrDatabase
	}

	condFunc := func(db *gorm.DB) *gorm.DB {
		return db.Where("block_number= ?", id)
	}

	txnList, err := uc.in.TransactionRepository.List(ctx, db, condFunc)
	if err != nil {
		uc.in.Logger.Error(ctx, err)
		return nil, codebook.ErrDatabase
	}

	return uc.arrangeResp(ctx, block, txnList)
}

func (uc *blockDetail) arrangeResp(ctx context.Context, block *model.Block, txnList []*model.Transaction) (*model.BlockDetail, error) {
	var (
		blockDetail model.BlockDetail
		txnHashList []string
	)

	if err := copier.Copy(&blockDetail, &block); err != nil {
		uc.in.Logger.Error(ctx, err)
		return nil, codebook.ErrServer
	}

	for _, v := range txnList {
		txnHashList = append(txnHashList, v.Hash)
	}

	blockDetail.Transactions = txnHashList

	return &blockDetail, nil
}

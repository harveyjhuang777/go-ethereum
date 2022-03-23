package block

import (
	"context"

	"github.com/jinzhu/copier"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/harveyjhuang777/go-ethereum/service/model"
	"github.com/harveyjhuang777/go-ethereum/service/util/codebook"
)

type ITransactionDetail interface {
	Handle(ctx context.Context, hash string) (*model.TransactionDetail, error)
}

func newTransactionDetail(in digIn) ITransactionDetail {
	return &transactionDetail{
		in: in,
	}
}

type transactionDetail struct {
	in digIn
}

func (uc *transactionDetail) Handle(ctx context.Context, hash string) (*model.TransactionDetail, error) {
	db := uc.in.DB.Session()

	transaction, err := uc.in.TransactionRepository.FirstByHash(ctx, db, hash)
	if err != nil {
		if xerrors.Is(gorm.ErrRecordNotFound, err) {
			return nil, codebook.ErrDataNotExist
		}
		uc.in.Logger.Error(ctx, err)
		return nil, codebook.ErrDatabase
	}

	condFunc := func(db *gorm.DB) *gorm.DB {
		return db.Where("transaction_hash = ?", transaction.Hash)
	}

	logs, err := uc.in.TransactionLogRepository.List(ctx, db, condFunc)
	if err != nil {
		uc.in.Logger.Error(ctx, err)
		return nil, codebook.ErrDatabase
	}

	return uc.arrangeResp(ctx, transaction, logs)
}

func (uc *transactionDetail) arrangeResp(ctx context.Context, transaction *model.Transaction, logs []*model.TransactionLog) (*model.TransactionDetail, error) {
	var (
		transactionDetail model.TransactionDetail
		logList           []*model.TransactionLogList
	)

	if err := copier.Copy(&transactionDetail, &transaction); err != nil {
		uc.in.Logger.Error(ctx, err)
		return nil, codebook.ErrServer
	}

	if err := copier.Copy(&logList, &logs); err != nil {
		uc.in.Logger.Error(ctx, err)
		return nil, codebook.ErrServer
	}

	transactionDetail.Logs = logList

	return &transactionDetail, nil
}

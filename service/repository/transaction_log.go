package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/harveyjhuang777/go-ethereum/service/model"
)

type ITransactionLog interface {
	Insert(ctx context.Context, db *gorm.DB, transactionLog *model.TransactionLog) error
	List(ctx context.Context, db *gorm.DB, condFunc ...func(*gorm.DB) *gorm.DB) ([]*model.TransactionLog, error)
}

type transactionLogRepository struct {
	in digIn
}

func newTransactionLogRepository(in digIn) ITransactionLog {
	return &transactionLogRepository{
		in: in,
	}
}

func (repo *transactionLogRepository) Insert(ctx context.Context, db *gorm.DB, transactionLog *model.TransactionLog) error {
	if err := db.Create(transactionLog).Error; err != nil {
		return err
	}
	return nil
}

func (repo *transactionLogRepository) List(ctx context.Context, db *gorm.DB, condFunc ...func(*gorm.DB) *gorm.DB) ([]*model.TransactionLog, error) {
	var resp []*model.TransactionLog

	if err := db.Scopes(condFunc...).Find(&resp).Error; err != nil {
		return nil, err
	}

	return resp, nil
}

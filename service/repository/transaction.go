package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/harveyjhuang777/go-ethereum/service/model"
)

type ITransaction interface {
	Insert(ctx context.Context, db *gorm.DB, transaction *model.Transaction) error
	Update(ctx context.Context, db *gorm.DB, transaction *model.Transaction) error
	First(ctx context.Context, db *gorm.DB, id int64) (*model.Transaction, error)
	List(ctx context.Context, db *gorm.DB, condFunc ...func(*gorm.DB) *gorm.DB) ([]*model.Transaction, error)
}

type transactionRepository struct {
	in digIn
}

func newTransactionRepository(in digIn) ITransaction {
	return &transactionRepository{
		in: in,
	}
}

func (repo *transactionRepository) Insert(ctx context.Context, db *gorm.DB, transaction *model.Transaction) error {
	if err := db.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (repo *transactionRepository) Update(ctx context.Context, db *gorm.DB, transaction *model.Transaction) error {
	if err := db.Save(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (repo *transactionRepository) First(ctx context.Context, db *gorm.DB, id int64) (*model.Transaction, error) {
	var resp model.Transaction
	if err := db.Where("id = ?", id).First(&resp).Error; err != nil {
		return nil, err
	}
	return &resp, nil
}

func (repo *transactionRepository) List(ctx context.Context, db *gorm.DB, condFunc ...func(*gorm.DB) *gorm.DB) ([]*model.Transaction, error) {
	var resp []*model.Transaction

	if err := db.Scopes(condFunc...).Find(&resp).Error; err != nil {
		return nil, err
	}

	return resp, nil
}

package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/harveyjhuang777/go-ethereum/service/model"
)

type IBlock interface {
	Insert(ctx context.Context, db *gorm.DB, block *model.Block) error
	Update(ctx context.Context, db *gorm.DB, block *model.Block) error
	First(ctx context.Context, db *gorm.DB, number int64) (*model.Block, error)
	List(ctx context.Context, db *gorm.DB, condFunc ...func(*gorm.DB) *gorm.DB) ([]*model.Block, error)
}

type blockRepository struct {
	in digIn
}

func newBlockRepository(in digIn) IBlock {
	return &blockRepository{
		in: in,
	}
}

func (repo *blockRepository) Insert(ctx context.Context, db *gorm.DB, block *model.Block) error {
	if err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(block).Error; err != nil {
		return err
	}
	return nil
}

func (repo *blockRepository) Update(ctx context.Context, db *gorm.DB, block *model.Block) error {
	if err := db.Save(block).Error; err != nil {
		return err
	}
	return nil
}

func (repo *blockRepository) First(ctx context.Context, db *gorm.DB, number int64) (*model.Block, error) {
	var resp model.Block
	if err := db.Where("number = ?", number).First(&resp).Error; err != nil {
		return nil, err
	}
	return &resp, nil
}

func (repo *blockRepository) List(ctx context.Context, db *gorm.DB, condFunc ...func(*gorm.DB) *gorm.DB) ([]*model.Block, error) {
	var resp []*model.Block

	if err := db.Scopes(condFunc...).Find(&resp).Error; err != nil {
		return nil, err
	}

	return resp, nil
}

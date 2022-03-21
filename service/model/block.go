package model

import (
	"time"
)

type Block struct {
	ID          int64         `gorm:"primaryKey"`
	Number      *string       `gorm:"column:number;type:varchar(10)"`
	Hash        *string       `gorm:"column:hash;type:varchar(70)"`
	Time        int64         `gorm:"column:time"`
	ParentHash  string        `gorm:"column:parent_hash;type:varchar(70)"`
	CreatedAt   time.Time     `gorm:"<-:create;column:created_at;autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	Transaction []Transaction `gorm:"save_associations:false;foreignKey:ID"`
}

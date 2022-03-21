package model

import "time"

type Transaction struct {
	ID        int64            `gorm:"primaryKey"`
	Hash      string           `gorm:"column:hash;type:varchar(70)"`
	From      string           `gorm:"column:from;type:varchar(50)"`
	To        string           `gorm:"column:to;type:varchar(50)"`
	Nonce     string           `gorm:"column:nonce;type:varchar(10)"`
	Data      string           `gorm:"column:data;type:varchar(10)"`
	Value     string           `gorm:"column:value;type:varchar(10)"`
	CreatedAt time.Time        `gorm:"<-:create;column:created_at;autoCreateTime"`
	UpdatedAt time.Time        `gorm:"column:updated_at;autoUpdateTime"`
	Logs      []TransactionLog `gorm:"save_associations:false;foreignKey:ID"`
}

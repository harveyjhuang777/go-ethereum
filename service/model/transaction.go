package model

import "time"

type Transaction struct {
	ID        int64            `gorm:"primaryKey"`
	BlockID   int64            `gorm:"column:block_id"`
	Hash      string           `gorm:"column:hash;type:varchar(70)"`
	From      string           `gorm:"column:from;type:varchar(50)"`
	To        string           `gorm:"column:to;type:varchar(50)"`
	Nonce     int64            `gorm:"column:nonce"`
	Data      string           `gorm:"column:data;type:varchar(70)"`
	Value     int64            `gorm:"column:value"`
	CreatedAt time.Time        `gorm:"<-:create;column:created_at;autoCreateTime"`
	UpdatedAt time.Time        `gorm:"column:updated_at;autoUpdateTime"`
	Logs      []TransactionLog `gorm:"save_associations:false;foreignKey:TransactionID"`
}

func (Transaction) TableName() string {
	return "transaction"
}

type TransactionDetail struct {
	Hash  string                `json:"tx_hash"`
	From  string                `json:"from"`
	To    string                `json:"to"`
	Nonce int64                 `json:"nonce"`
	Data  string                `json:"data"`
	Value int64                 `json:"value"`
	Logs  []*TransactionLogList `json:"logs"`
}

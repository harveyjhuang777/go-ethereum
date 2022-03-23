package model

import (
	"time"
)

type Block struct {
	Number      uint64        `gorm:"primaryKey;column:number"`
	Hash        string        `gorm:"column:hash;type:varchar(70)"`
	Time        uint64        `gorm:"column:time"`
	ParentHash  string        `gorm:"column:parent_hash;type:varchar(70)"`
	CreatedAt   time.Time     `gorm:"<-:create;column:created_at;autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	Transaction []Transaction `gorm:"save_associations:false;foreignKey:BlockNumber"`
}

func (Block) TableName() string {
	return "block"
}

type BlockListResp struct {
	Blocks []*BlockList `json:"blocks"`
}

type BlockList struct {
	Number     uint64 `json:"block_number"`
	Hash       string `json:"block_hash"`
	Time       int64  `json:"block_time"`
	ParentHash string `json:"parent_hash"`
}

type BlockDetail struct {
	Number       uint64   `json:"block_number"`
	Hash         string   `json:"block_hash"`
	Time         int64    `json:"block_time"`
	ParentHash   string   `json:"parent_hash"`
	Transactions []string `json:"transactions"`
}

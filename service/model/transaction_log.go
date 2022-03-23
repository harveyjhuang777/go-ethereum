package model

import "time"

type TransactionLog struct {
	ID              int64     `gorm:"primaryKey"`
	TransactionHash string    `gorm:"transaction_hash;type:varchar(70)"`
	Index           uint      `gorm:"column:index"`
	Data            string    `gorm:"column:data;type:TEXT"`
	CreatedAt       time.Time `gorm:"<-:create;column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (TransactionLog) TableName() string {
	return "transaction_log"
}

type TransactionLogList struct {
	Index int64  `json:"index"`
	Data  string `json:"data"`
}

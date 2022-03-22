package model

import "time"

type TransactionLog struct {
	ID            int64     `gorm:"primaryKey"`
	TransactionID int64     `gorm:"transaction_id"`
	Index         int64     `gorm:"column:index"`
	Data          string    `gorm:"column:data;type:varchar(70)"`
	CreatedAt     time.Time `gorm:"<-:create;column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (TransactionLog) TableName() string {
	return "transaction_log"
}

package model

import "time"

type TransactionLog struct {
	ID            int64     `gorm:"primaryKey"`
	TransactionID int64     `gorm:"transaction_id"`
	Index         string    `gorm:"column:index;type:varchar(10)"`
	Data          string    `gorm:"column:data;type:varchar(70)"`
	CreatedAt     time.Time `gorm:"<-:create;column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

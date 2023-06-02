package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID     uint
	OfficeID   uint
	Discount   int
	Price      int
	Tax        int
	TotalPrice int
	Start      time.Time
	End        time.Time
	Status     bool
}

type TransactionParam struct {
	ID       uint `uri:"transaction_id" json:"id"`
	UserID   uint
	OfficeID uint
	Start    time.Time
	End      time.Time
}

type UpdateTransactionParam struct {
	Start  time.Time
	End    time.Time
	Status bool
}

type CreateTransactionParam struct {
	OfficeID  uint   `json:"-" uri:"office_id"`
	Start     string `binding:"required"`
	End       string `binding:"required"`
	PaymentID string `binding:"required"`
}

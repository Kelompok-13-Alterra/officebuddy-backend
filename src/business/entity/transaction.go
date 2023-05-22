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
	Tax        int
	TotalPrice int
	Start      time.Time
	End        time.Time
	Status     bool
}

type TransactionParam struct {
	ID       uint
	UserID   uint
	OfficeID uint
}

type UpdateTransactionParam struct {
	Start  time.Time
	End    time.Time
	Status bool
}

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

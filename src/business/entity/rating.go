package entity

import (
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	UserID      uint
	OfficeID    uint
	Star        int
	Tags        string
	Description string `gorm:"type:text"`
}

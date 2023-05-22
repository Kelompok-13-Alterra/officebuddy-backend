package entity

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	OfficeID uint
	UserID   uint
	Step     string
	Status   bool
}

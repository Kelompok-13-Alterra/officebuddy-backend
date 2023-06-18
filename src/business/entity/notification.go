package entity

import (
	"gorm.io/gorm"
)

const (
	WelcomeStatus    = "welcome"
	FirstOrderStatus = "firstorder"
	ProcessingStatus = "processing"
	SuccessStatus    = "success"
)

type Notification struct {
	gorm.Model
	UserID      uint
	Description string
	Status      string
	IsRead      bool
}

type NotificationParam struct {
	ID       uint
	UserID   uint
	OfficeID uint
	OrderBy  string `gorm:"-"`
}

type UpdateNotificationParam struct {
	IsRead	bool
}

package entity

import (
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	UserID        uint
	OfficeID      uint
	TransactionID uint
	Star          int
	Tags          string
	Description   string
}

type RatingParam struct {
	ID            uint `json:"-" uri:"id" binding:"required"`
	UserID        uint
	OfficeID      uint
	TransactionID uint
}

type CreateRatingParam struct {
	TransactionID uint     `json:"-" uri:"transaction_id"`
	Star          int      `binding:"required"`
	Tags          []string `binding:"required"`
	Description   string   `binding:"required"`
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	UserID        uint
	OfficeID      uint
	TransactionID uint
	Star          float32
	Tags          string
	Description   string
}

type RatingResponse struct {
	ID            uint
	UserID        uint
	OfficeID      uint
	TransactionID uint
	Star          float32
	Tags          []string
	Description   string
	CreatedAt     time.Time
}

type RatingParam struct {
	ID            uint `json:"-" uri:"rating_id"`
	UserID        uint
	OfficeID      uint
	TransactionID uint
}

type CreateRatingParam struct {
	TransactionID uint     `json:"-" uri:"transaction_id"`
	Star          float32  `binding:"required"`
	Tags          []string `binding:"required"`
	Description   string   `binding:"required"`
}

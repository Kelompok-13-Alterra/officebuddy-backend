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
	Description string
}

type RatingParam struct {
	ID       uint
	UserID   uint
	OfficeID uint
}

type CreateRatingParam struct {
	UserID      uint
	OfficeID    uint   `json:"-" uri:"office_id"`
	Star        int    `binding:"required"`
	Tags        string `binding:"required"`
	Description string `binding:"required"`
}

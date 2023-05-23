package entity

import (
	"gorm.io/gorm"
)

type ImageRating struct {
	gorm.Model
	UserID   uint
	RatingID uint
	Url      string
}

type ImageRatingParam struct {
	ID       uint
	UserID   uint
	RatingID uint
}

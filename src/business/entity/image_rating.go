package entity

import (
	"gorm.io/gorm"
)

type ImageRating struct {
	gorm.Model
	UserID		uint
	RatingID	uint
	url			string
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	Name        string
	Description string `gorm:"type:text"`
	Capacity    int
	Open        time.Time
	Close       time.Time
	Price       int
	Location    string
	Facilities  string
	Status      bool
}

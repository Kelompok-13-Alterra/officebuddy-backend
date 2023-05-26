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
	Open        time.Time `gorm:"type:time"`
	Close       time.Time `gorm:"type:time"`
	Price       int
	Location    string
	Facilities  string
	Status      bool
}
type CreateOfficeParam struct {
	Name        string
	Description string
	Location    string
	Facilities  string
}

type OfficeParam struct {
	ID   uint   `uri:"office_id" json:"id"`
	Name string `form:"name" json:"name"`
}

type UpdateOfficeParam struct {
	Open   time.Time
	Close  time.Time
	Status bool
}

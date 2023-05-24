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
type CreateOfficeParam struct {
	Name        string
	Description string
	Location    string
	Facilities  string
}

type OfficeParam struct {
	Name		string
	Location	string
	Facilities	string
	Price		int
	Status		bool
}

type UpdateOfficeParam struct {
	Open	time.Time
	Close	time.Time
	Status	bool
}

package entity

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	Name        string
	Description string `gorm:"type:text"`
	Capacity    int
	Type        string
	Open        OfficeHours
	Close       OfficeHours
	Price       int
	Location    string
	Facilities  string
	Status      bool
}

type OfficeHours struct {
	time.Time
}

func (t *OfficeHours) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.Time.UnmarshalText(v)
	case string:
		return t.Time.UnmarshalText([]byte(v))
	default:
		return errors.New("type not supported")
	}
}

func (t OfficeHours) Value() (driver.Value, error) {
	return t.Format("15:04:05"), nil
}

type CreateOfficeParam struct {
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Capacity    int    `binding:"required"`
	Type        string `json:"-"`
	Open        string `binding:"required"`
	Close       string `binding:"required"`
	Location    string `binding:"required"`
	Price       int    `binding:"required"`
	Facilities  string `binding:"required"`
}

type OfficeTypeParam struct {
	Type string `form:"type" binding:"required"`
}

type OfficeParam struct {
	ID   uint   `uri:"office_id" json:"id"`
	Name string `form:"name" json:"name"`
}

type UpdateOfficeParam struct {
	Open   string
	Close  string
	Status bool
}

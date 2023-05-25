package domain

import (
	"go-clean/src/business/domain/office"
	"go-clean/src/business/domain/user"

	"gorm.io/gorm"
)

type Domains struct {
	User   user.Interface
	Office office.Interface
}

func Init(db *gorm.DB) *Domains {
	d := &Domains{
		User:   user.Init(db),
		Office: office.Init(db),
	}

	return d
}

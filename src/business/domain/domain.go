package domain

import (
	"go-clean/src/business/domain/notification"
	"go-clean/src/business/domain/office"
	"go-clean/src/business/domain/transaction"
	"go-clean/src/business/domain/user"

	"gorm.io/gorm"
)

type Domains struct {
	User         user.Interface
	Office       office.Interface
	Transaction  transaction.Interface
	Notification notification.Interface
}

func Init(db *gorm.DB) *Domains {
	d := &Domains{
		User:         user.Init(db),
		Office:       office.Init(db),
		Transaction:  transaction.Init(db),
		Notification: notification.Init(db),
	}

	return d
}

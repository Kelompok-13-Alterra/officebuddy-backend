package domain

import (
	"go-clean/src/business/domain/image_rating"
	"go-clean/src/business/domain/midtrans_transaction"
	"go-clean/src/business/domain/notification"
	"go-clean/src/business/domain/office"
	"go-clean/src/business/domain/rating"
	"go-clean/src/business/domain/transaction"
	"go-clean/src/business/domain/user"

	"gorm.io/gorm"
)

type Domains struct {
	User                user.Interface
	Office              office.Interface
	Transaction         transaction.Interface
	Notification        notification.Interface
	Rating              rating.Interface
	MidtransTransaction midtrans_transaction.Interface
	ImageRating         image_rating.Interface
}

func Init(db *gorm.DB) *Domains {
	d := &Domains{
		User:                user.Init(db),
		Office:              office.Init(db),
		Transaction:         transaction.Init(db),
		Notification:        notification.Init(db),
		Rating:              rating.Init(db),
		MidtransTransaction: midtrans_transaction.Init(db),
		ImageRating:         image_rating.Init(db),
	}

	return d
}

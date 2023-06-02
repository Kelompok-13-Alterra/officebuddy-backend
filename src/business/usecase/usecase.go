package usecase

import (
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase/midtrans_transaction"
	"go-clean/src/business/usecase/notification"
	"go-clean/src/business/usecase/office"
	"go-clean/src/business/usecase/rating"
	"go-clean/src/business/usecase/transaction"
	"go-clean/src/business/usecase/user"
	"go-clean/src/lib/auth"
)

type Usecase struct {
	User                user.Interface
	Office              office.Interface
	Transaction         transaction.Interface
	Notification        notification.Interface
	MidtransTransaction midtrans_transaction.Interface
	Rating              rating.Interface
}

func Init(auth auth.Interface, d *domain.Domains) *Usecase {
	uc := &Usecase{
		User:                user.Init(d.User, auth),
		Office:              office.Init(d.Office),
		Transaction:         transaction.Init(d.Transaction, auth, d.Office, d.Midtrans, d.MidtransTransaction),
		Notification:        notification.Init(d.Notification, auth),
		MidtransTransaction: midtrans_transaction.Init(d.MidtransTransaction, d.Midtrans),
		Rating:              rating.Init(d.Rating),
	}

	return uc
}

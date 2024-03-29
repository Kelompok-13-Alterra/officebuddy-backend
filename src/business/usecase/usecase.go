package usecase

import (
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase/midtrans_transaction"
	"go-clean/src/business/usecase/notification"
	"go-clean/src/business/usecase/office"
	"go-clean/src/business/usecase/rating"
	"go-clean/src/business/usecase/transaction"
	"go-clean/src/business/usecase/user"
	"go-clean/src/business/usecase/widget_analytic"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/cloud_storage"
)

type Usecase struct {
	User                user.Interface
	Office              office.Interface
	Rating              rating.Interface
	Transaction         transaction.Interface
	Notification        notification.Interface
	MidtransTransaction midtrans_transaction.Interface
	WidgetDashboard     widget_analytic.Interface
}

func Init(auth auth.Interface, d *domain.Domains, cs cloud_storage.Interface) *Usecase {
	uc := &Usecase{
		User:                user.Init(d.User, auth, d.Notification),
		Office:              office.Init(d.Office, cs),
		Transaction:         transaction.Init(d.Transaction, auth, d.Office, d.Midtrans, d.MidtransTransaction, d.User),
		Notification:        notification.Init(d.Notification, auth),
		Rating:              rating.Init(d.Rating, d.Transaction, auth),
		MidtransTransaction: midtrans_transaction.Init(d.Transaction, d.Office, d.MidtransTransaction, d.Midtrans, d.Notification),
		WidgetDashboard:     widget_analytic.Init(d.Office, d.Transaction, d.Rating, d.MidtransTransaction, d.User),
	}

	return uc
}

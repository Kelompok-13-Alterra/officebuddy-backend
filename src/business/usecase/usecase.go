package usecase

import (
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase/notification"
	"go-clean/src/business/usecase/office"
	"go-clean/src/business/usecase/transaction"
	"go-clean/src/business/usecase/user"
	"go-clean/src/lib/auth"
)

type Usecase struct {
	User        user.Interface
	Office      office.Interface
	Transaction transaction.Interface
  Notification notification.Interface
}

func Init(auth auth.Interface, d *domain.Domains) *Usecase {
	uc := &Usecase{
		User:        user.Init(d.User, auth),
		Office:      office.Init(d.Office),
		Transaction: transaction.Init(d.Transaction, auth),
    Notification: notification.Init(d.Notification, auth),
	}

	return uc
}

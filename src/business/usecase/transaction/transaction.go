package transaction

import (
	"context"
	transactionDom "go-clean/src/business/domain/transaction"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
)

type Interface interface {
	GetListBooked(ctx context.Context, param entity.TransactionParam) ([]entity.Transaction, error)
}

type transaction struct {
	transaction transactionDom.Interface
	auth        auth.Interface
}

func Init(od transactionDom.Interface, auth auth.Interface) Interface {
	t := &transaction{
		transaction: od,
		auth:        auth,
	}

	return t
}

func (t *transaction) GetListBooked(ctx context.Context, param entity.TransactionParam) ([]entity.Transaction, error) {
	var (
		transactions []entity.Transaction
		err          error
	)

	user, err := t.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return transactions, err
	}

	transactions, err = t.transaction.GetListBooked(entity.TransactionParam{
		UserID: user.User.ID,
	})
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

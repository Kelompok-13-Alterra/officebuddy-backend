package transaction

import (
	"context"
	"encoding/json"
	"errors"
	midtransDom "go-clean/src/business/domain/midtrans"
	midtransTransactionDom "go-clean/src/business/domain/midtrans_transaction"
	officeDom "go-clean/src/business/domain/office"
	transactionDom "go-clean/src/business/domain/transaction"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/midtrans"
	"strconv"
	"time"

	"github.com/midtrans/midtrans-go/coreapi"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateTransactionParam) (uint, error)
	GetListBooked(ctx context.Context) ([]entity.Transaction, error)
}

type transaction struct {
	transaction         transactionDom.Interface
	auth                auth.Interface
	office              officeDom.Interface
	midtrans            midtransDom.Interface
	midtransTransaction midtransTransactionDom.Interface
}

func Init(td transactionDom.Interface, auth auth.Interface, od officeDom.Interface, md midtransDom.Interface, mtd midtransTransactionDom.Interface) Interface {
	t := &transaction{
		transaction:         td,
		auth:                auth,
		office:              od,
		midtrans:            md,
		midtransTransaction: mtd,
	}

	return t
}

func (t *transaction) Create(ctx context.Context, param entity.CreateTransactionParam) (uint, error) {
	user, err := t.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return 0, err
	}

	office, err := t.office.Get(entity.OfficeParam{
		ID: param.OfficeID,
	})
	if err != nil {
		return 0, err
	}

	startDate, endDate, err := t.formatDate(param.Start, param.End)
	if err != nil {
		return 0, err
	}

	transactionExist, err := t.transaction.GetAvaibility(entity.TransactionParam{
		Start:    startDate,
		End:      endDate,
		OfficeID: param.OfficeID,
	})
	if err != nil {
		return 0, err
	}

	if transactionExist.ID != 0 {
		return 0, errors.New("tanggal tidak tersedia")
	}

	diff := endDate.Sub(startDate)
	days := int(diff.Hours()/24) + 1

	officePrice := office.Price * days
	tax := officePrice * 10 / 100
	totalPrice := officePrice + tax

	transaction, err := t.transaction.Create(entity.Transaction{
		UserID:     user.User.ID,
		OfficeID:   office.ID,
		Start:      startDate,
		End:        endDate,
		Price:      officePrice,
		Discount:   0,
		Tax:        tax,
		TotalPrice: totalPrice,
	})
	if err != nil {
		return 0, err
	}

	coreApiRes, err := t.midtrans.Create(midtrans.CreateOrderParam{
		OrderID:     transaction.ID,
		PaymentID:   param.PaymentID,
		GrossAmount: int64(totalPrice),
		ItemsDetails: midtrans.ItemsDetails{
			ID:    strconv.Itoa(int(office.ID)),
			Price: int64(office.Price) * 110 / 100,
			Qty:   days,
			Name:  office.Name,
		},
		CustomerDetails: midtrans.CustomerDetails{
			Name:  user.User.Name,
			Email: user.User.Email,
		},
	})
	if err != nil {
		return 0, err
	}

	paymentData, err := t.getPaymentData(param.PaymentID, coreApiRes, totalPrice, officePrice, 0, tax)
	if err != nil {
		return 0, err
	}

	paymentDataMarshal, err := json.Marshal(paymentData)
	if err != nil {
		return 0, err
	}

	_, err = t.midtransTransaction.Create(entity.MidtransTransaction{
		TransactionID: transaction.ID,
		MidtransID:    coreApiRes.TransactionID,
		OrderID:       coreApiRes.OrderID,
		PaymentType:   param.PaymentID,
		Amount:        totalPrice,
		Status:        entity.StatusPending,
		PaymentData:   string(paymentDataMarshal),
	})
	if err != nil {
		return 0, err
	}

	return transaction.ID, nil
}

func (t *transaction) getPaymentData(paymentId string, coreApiRes *coreapi.ChargeResponse, totalPrice, price, discount, tax int) (entity.PaymentData, error) {
	paymentData := entity.PaymentData{
		TotalPrice: totalPrice,
		Discount:   discount,
		Tax:        tax,
		Price:      price,
	}
	if paymentId == midtrans.VaBNI {
		paymentData.VaNumber = coreApiRes.VaNumbers[0].VANumber
	} else {
		return paymentData, errors.New("failed to get payment data")
	}

	return paymentData, nil
}

func (t *transaction) formatDate(start, end string) (time.Time, time.Time, error) {
	var startDate, endDate time.Time

	layoutFormat := "2006-01-02"
	startDate, err := time.Parse(layoutFormat, start)
	if err != nil {
		return startDate, endDate, err
	}

	endDate, err = time.Parse(layoutFormat, end)
	if err != nil {
		return startDate, endDate, err
	}

	return startDate, endDate, nil
}

func (t *transaction) GetListBooked(ctx context.Context) ([]entity.Transaction, error) {
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

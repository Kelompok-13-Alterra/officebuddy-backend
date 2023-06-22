package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	midtransDom "go-clean/src/business/domain/midtrans"
	midtransTransactionDom "go-clean/src/business/domain/midtrans_transaction"
	officeDom "go-clean/src/business/domain/office"
	transactionDom "go-clean/src/business/domain/transaction"
	userDom "go-clean/src/business/domain/user"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/midtrans"
	"go-clean/src/lib/timeutils"
	"log"
	"strconv"
	"time"

	"github.com/midtrans/midtrans-go/coreapi"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateTransactionParam) (uint, error)
	GetTransactionList(param entity.TransactionParam) ([]entity.Transaction, error)
	GetListBooked(ctx context.Context) ([]entity.Transaction, error)
	GetListHistoryBooked(ctx context.Context) ([]entity.Transaction, error)
	RescheduleBooked(ctx context.Context, param entity.InputUpdateTransactionParam, selectParam entity.TransactionParam) error
	ValidateTransaction(ctx context.Context, transactionID uint, userID uint) error
	AvailabilityCheck(param entity.AvailabilityCheckTransactionParam) error
	GetLastTransactionList(param entity.MidtransTransactionParam) ([]entity.LastTranasctionResult, error)
}

type transaction struct {
	transaction         transactionDom.Interface
	auth                auth.Interface
	office              officeDom.Interface
	midtrans            midtransDom.Interface
	midtransTransaction midtransTransactionDom.Interface
	user                userDom.Interface
}

func Init(td transactionDom.Interface, auth auth.Interface, od officeDom.Interface, md midtransDom.Interface, mtd midtransTransactionDom.Interface, ud userDom.Interface) Interface {
	t := &transaction{
		transaction:         td,
		auth:                auth,
		office:              od,
		midtrans:            md,
		midtransTransaction: mtd,
		user:                ud,
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

func (t *transaction) AvailabilityCheck(param entity.AvailabilityCheckTransactionParam) error {
	startDate, endDate, err := t.formatDate(param.Start, param.End)
	if err != nil {
		return err
	}

	transactionExist, err := t.transaction.GetAvaibility(entity.TransactionParam{
		Start:    startDate,
		End:      endDate,
		OfficeID: param.OfficeID,
	})
	if err != nil {
		return err
	}

	if transactionExist.ID != 0 {
		return errors.New("tanggal tidak tersedia")
	}

	return nil
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

func (t *transaction) GetTransactionList(param entity.TransactionParam) ([]entity.Transaction, error) {
	var (
		transactions []entity.Transaction
		err          error
	)

	transactions, err = t.transaction.GetList(entity.TransactionParam{})

	for i, trans := range transactions {
		midtransTransactions, _ := t.midtransTransaction.Get(entity.MidtransTransactionParam{TransactionID: trans.ID})
		transactions[i].PaymentStatus = midtransTransactions.Status
	}

	if err != nil {
		return transactions, err
	}

	return transactions, nil
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

	mapOfficeIDs := make(map[uint]bool)
	for _, t := range transactions {
		mapOfficeIDs[t.OfficeID] = true
	}

	officeIDs := []uint{}
	for id := range mapOfficeIDs {
		officeIDs = append(officeIDs, id)
	}

	offices, err := t.office.GetListByID(officeIDs)
	if err != nil {
		return transactions, err
	}

	officesMap := make(map[uint]entity.Office)
	for _, o := range offices {
		if o.ImageUrl != "" {
			url, err := t.office.GetPresignedURL(ctx, o.ImageUrl)
			if err != nil {
				log.Println(err)
			}
			o.ImageUrl = url
		}
		officesMap[o.ID] = o
	}

	for i, t := range transactions {
		transactions[i].Office = officesMap[t.OfficeID]
	}

	return transactions, nil
}

func (t *transaction) GetListHistoryBooked(ctx context.Context) ([]entity.Transaction, error) {
	var (
		transactions []entity.Transaction
		err          error
	)

	user, err := t.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return transactions, err
	}

	transactions, err = t.transaction.GetListHistoryBooked(entity.TransactionParam{
		UserID: user.User.ID,
	})
	if err != nil {
		return transactions, err
	}

	mapOfficeIDs := make(map[uint]bool)
	for _, t := range transactions {
		mapOfficeIDs[t.OfficeID] = true
	}

	officeIDs := []uint{}
	for id := range mapOfficeIDs {
		officeIDs = append(officeIDs, id)
	}

	offices, err := t.office.GetListByID(officeIDs)
	if err != nil {
		return transactions, err
	}

	officesMap := make(map[uint]entity.Office)
	for _, o := range offices {
		if o.ImageUrl != "" {
			url, err := t.office.GetPresignedURL(ctx, o.ImageUrl)
			if err != nil {
				log.Println(err)
			}
			o.ImageUrl = url
		}
		officesMap[o.ID] = o
	}

	for i, t := range transactions {
		transactions[i].Office = officesMap[t.OfficeID]
	}

	return transactions, nil
}

func (t *transaction) RescheduleBooked(ctx context.Context, param entity.InputUpdateTransactionParam, selectParam entity.TransactionParam) error {
	startDate, endDate, err := t.formatDate(param.Start, param.End)
	if err != nil {
		return err
	}

	transaction, err := t.transaction.Get(entity.TransactionParam{
		ID: selectParam.ID,
	})
	if err != nil {
		return err
	}

	transactionExist, err := t.transaction.GetAvaibility(entity.TransactionParam{
		Start:    startDate,
		End:      endDate,
		OfficeID: transaction.OfficeID,
	})
	if err != nil {
		return err
	}

	if transactionExist.ID != 0 {
		return errors.New("tanggal tidak tersedia")
	}

	if err := t.transaction.Update(entity.TransactionParam{
		ID: selectParam.ID,
	}, entity.UpdateTransactionParam{
		Start: startDate,
		End:   endDate,
	}); err != nil {
		return err
	}

	return nil
}

func (t *transaction) ValidateTransaction(ctx context.Context, transactionID uint, userID uint) error {
	transaction, err := t.transaction.Get(entity.TransactionParam{
		ID: transactionID,
	})
	if err != nil {
		return err
	}

	if transaction.UserID != userID {
		return errors.New("unauthorized")
	}

	return nil
}

func (t *transaction) GetLastTransactionList(param entity.MidtransTransactionParam) ([]entity.LastTranasctionResult, error) {
	result := []entity.LastTranasctionResult{}

	mTransactions, err := t.midtransTransaction.GetListWithPaginationByIDs(entity.MidtransTransactionParam{
		Status:  entity.StatusSuccess,
		Limit:   param.Limit,
		Offset:  (param.Page - 1) * param.Limit,
		OrderBy: "id desc",
	})
	if err != nil {
		return result, err
	}

	transactionIDs := []uint{}
	for _, mt := range mTransactions {
		transactionIDs = append(transactionIDs, mt.TransactionID)
	}

	transactions, err := t.transaction.GetListByIDs(transactionIDs)
	if err != nil {
		return result, err
	}

	userIdsMap := make(map[uint]bool)
	transactionMap := make(map[uint]entity.Transaction)
	officeIDsMap := make(map[uint]bool)
	for _, t := range transactions {
		transactionMap[t.ID] = t
		userIdsMap[t.UserID] = true
		officeIDsMap[t.OfficeID] = true
	}

	userIds := []uint{}
	for id := range userIdsMap {
		userIds = append(userIds, id)
	}

	users, err := t.user.GetListByIDs(userIds)
	if err != nil {
		return result, err
	}

	userMap := make(map[uint]entity.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	officeIDs := []uint{}
	for id := range officeIDsMap {
		officeIDs = append(officeIDs, id)
	}

	offices, err := t.office.GetListByID(officeIDs)
	if err != nil {
		return result, err
	}

	officesMap := make(map[uint]entity.Office)
	for _, o := range offices {
		officesMap[o.ID] = o
	}

	for _, mt := range mTransactions {
		officeType := ""
		office := officesMap[transactionMap[mt.TransactionID].OfficeID]
		if office.Type == entity.CoWorkingType {
			officeType = "Coworking Space"
		} else {
			officeType = "Office"
		}
		result = append(result, entity.LastTranasctionResult{
			BuyerName:   userMap[transactionMap[mt.TransactionID].UserID].Name,
			Description: fmt.Sprintf("Memesan %s %s", officeType, office.Name),
			Revenue:     mt.Amount,
			Date:        timeutils.DiffForHumans(mt.CreatedAt),
		})
	}

	return result, nil
}

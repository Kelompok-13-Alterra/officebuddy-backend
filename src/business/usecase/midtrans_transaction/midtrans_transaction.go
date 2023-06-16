package midtrans_transaction

import (
	"encoding/json"
	"errors"
	"fmt"
	midtransDom "go-clean/src/business/domain/midtrans"
	midtransTransactionDom "go-clean/src/business/domain/midtrans_transaction"
	notificationDom "go-clean/src/business/domain/notification"
	transactionDom "go-clean/src/business/domain/transaction"
	"go-clean/src/business/entity"
	"log"
)

type Interface interface {
	GetPaymentDetail(param entity.MidtransTransactionParam) (entity.MidtransTransactionPaymentDetail, error)
	HandleNotification(payload map[string]interface{}) error
}

type midtransTransaction struct {
	midtransTransaction midtransTransactionDom.Interface
	midtrans            midtransDom.Interface
	notification        notificationDom.Interface
	transaction         transactionDom.Interface
}

func Init(mttd midtransTransactionDom.Interface, md midtransDom.Interface, td transactionDom.Interface) Interface {
	mtt := &midtransTransaction{
		midtransTransaction: mttd,
		midtrans:            md,
		transaction:         td,
	}

	return mtt
}

func (mtt *midtransTransaction) GetPaymentDetail(param entity.MidtransTransactionParam) (entity.MidtransTransactionPaymentDetail, error) {
	result := entity.MidtransTransactionPaymentDetail{}

	midtransTransaction, err := mtt.midtransTransaction.Get(param)
	if err != nil {
		return result, err
	}

	paymentData := entity.PaymentData{}
	if err := json.Unmarshal([]byte(midtransTransaction.PaymentData), &paymentData); err != nil {
		return result, err
	}

	result.Status = midtransTransaction.Status
	result.PaymentData = paymentData
	result.PaymentType = midtransTransaction.PaymentType

	return result, nil
}

func (mtt *midtransTransaction) HandleNotification(payload map[string]interface{}) error {
	orderId, exist := payload["order_id"].(string)
	if !exist {
		return errors.New("order id not exist")
	}

	transactionResponse, err := mtt.midtrans.HandleNotification(orderId)
	if err != nil {
		return err
	}

	status := ""

	if transactionResponse != nil {
		// 5. Do set transaction status based on response from check transaction status
		if transactionResponse.TransactionStatus == "capture" {
			if transactionResponse.FraudStatus == "challenge" {
				// TODO set transaction status on your database to 'challenge'
				status = entity.StatusChallange
				// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			} else if transactionResponse.FraudStatus == "accept" {
				// TODO set transaction status on your database to 'success'
				status = entity.StatusSuccess
			}
		} else if transactionResponse.TransactionStatus == "settlement" {
			// TODO set transaction status on your databaase to 'success'
			status = entity.StatusSuccess
		} else if transactionResponse.TransactionStatus == "deny" {
			// TODO you can ignore 'deny', because most of the time it allows payment retries
			// and later can become success
			status = entity.StatusDeny
		} else if transactionResponse.TransactionStatus == "cancel" || transactionResponse.TransactionStatus == "expire" {
			// TODO set transaction status on your databaase to 'failure'
			status = entity.StatusFailure
		} else if transactionResponse.TransactionStatus == "pending" {
			// TODO set transaction status on your databaase to 'pending' / waiting payment
			status = entity.StatusPending
		}
	}

	if err := mtt.midtransTransaction.Update(entity.MidtransTransactionParam{
		OrderID: orderId,
	}, entity.UpdateMidtransTransactionParam{
		Status: status,
	}); err != nil {
		return err
	}

	var transactionId uint
	_, err = fmt.Sscanf(orderId, "OB-%d-00000000", &transactionId)
	if err != nil {
		log.Println("failed to scan transaction id")
	}

	transaction, err := mtt.transaction.Get(entity.TransactionParam{
		ID: transactionId,
	})
	if err != nil {
		log.Println("failed to get transaction id")
	} else {
		_, err := mtt.notification.Create(entity.Notification{
			UserID:      transaction.UserID,
			Description: "Bookingan kamu <b>sedang diproses</b>",
			Status:      entity.ProcessingStatus,
			IsRead:      false,
		})
		if err != nil {
			log.Println("failed to create processing notification")
		}
		_, err = mtt.notification.Create(entity.Notification{
			UserID:      transaction.UserID,
			Description: fmt.Sprintf("Booking Office dengan No Pesanan <b>#%d</b> Berhasil", transaction.ID),
			Status:      entity.SuccessStatus,
			IsRead:      false,
		})
		if err != nil {
			log.Println("failed to create processing notification")
		}
	}

	return nil
}

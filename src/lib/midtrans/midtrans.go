package midtrans

import (
	"errors"
	"fmt"
	"time"

	midtransSdk "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type Interface interface {
	CreateOrder(param CreateOrderParam) (*coreapi.ChargeResponse, error)
	HandleNotification(id string) (*coreapi.TransactionStatusResponse, error)
}

type Config struct {
	ServerKey string
}

type midtrans struct {
	conf    Config
	coreapi *coreapi.Client
}

func Init(cfg Config) Interface {
	m := &midtrans{
		conf: cfg,
	}
	m.connect()
	return m
}

func (m *midtrans) connect() {
	c := coreapi.Client{}
	c.New(m.conf.ServerKey, midtransSdk.Sandbox)
	m.coreapi = &c
}

func (m *midtrans) CreateOrder(param CreateOrderParam) (*coreapi.ChargeResponse, error) {
	fmt.Printf("isi param : %#v", param)
	chargeReq := &coreapi.ChargeReq{
		TransactionDetails: midtransSdk.TransactionDetails{
			OrderID:  fmt.Sprintf("%s-%d-%d", "OB", param.OrderID, time.Now().Unix()),
			GrossAmt: param.GrossAmount,
		},
		Items: &[]midtransSdk.ItemDetails{
			{
				ID:    param.ItemsDetails.ID,
				Name:  param.ItemsDetails.Name,
				Price: param.ItemsDetails.Price,
				Qty:   int32(param.ItemsDetails.Qty),
			},
		},
		CustomerDetails: &midtransSdk.CustomerDetails{
			FName: param.CustomerDetails.Name,
			Email: param.CustomerDetails.Email,
		},
	}

	if param.PaymentID == VaBNI {
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtransSdk.BankBni,
		}
	} else {
		return &coreapi.ChargeResponse{}, errors.New("undeifned payment method")
	}

	coreApiRes, err := m.coreapi.ChargeTransaction(chargeReq)
	if err != nil {
		return coreApiRes, err
	}

	return coreApiRes, nil
}

func (m *midtrans) HandleNotification(id string) (*coreapi.TransactionStatusResponse, error) {
	midtransReport, err := m.coreapi.CheckTransaction(id)
	if err != nil {
		return midtransReport, err
	}

	return midtransReport, nil
}

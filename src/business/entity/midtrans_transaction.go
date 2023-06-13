package entity

import "gorm.io/gorm"

const (
	StatusChallange = "challange"
	StatusSuccess   = "success"
	StatusDeny      = "deny"
	StatusFailure   = "failure"
	StatusPending   = "pending"
)

type MidtransTransaction struct {
	gorm.Model
	TransactionID uint
	MidtransID    string
	OrderID       string
	PaymentType   string
	Amount        int
	Status        string
	PaymentData   string
}

type PaymentData struct {
	Key        string `json:"key,omitempty"`
	Qr         string `json:"qr,omitempty"`
	VaNumber   string `json:"va_number,omitempty"`
	TotalPrice int    `json:"total_price"`
	Discount   int    `json:"discount"`
	Tax        int    `json:"tax"`
	Price      int    `json:"price"`
}

type MidtransTransactionPaymentDetail struct {
	PaymentType string      `json:"payment_type"`
	PaymentData PaymentData `json:"payment_data"`
	Status      string      `json:"status"`
}

type MidtransTransactionParam struct {
	TransactionID uint   `uri:"transaction_id" json:"transaction_id"`
	OrderID       string `json:"order_id"`
	Status        string
	Limit         int    `form:"limit" gorm:"-" json:"-"`
	Page          int    `form:"page" gorm:"-" json:"-"`
	OrderBy       string `gorm:"-" json:"-"`
	Offset        int    `gorm:"-" json:"-"`
}

type UpdateMidtransTransactionParam struct {
	Status string `json:"string"`
}

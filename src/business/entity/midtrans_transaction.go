package entity

import "gorm.io/gorm"

type MidtransTransaction struct {
	gorm.Model
	TransactionID uint
	MidtransID    string
	PaymentType   string
	Amount        int
	Status        string
	PaymentData   string
}

type MidtransTransactionParam struct {
	ID            uint
	TransactionID uint
	MidtransID    string
}

type UpdateMidtransTransactionParam struct {
	Status string
}

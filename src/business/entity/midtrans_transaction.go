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

package midtrans_transaction

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(transaction entity.MidtransTransaction) (entity.MidtransTransaction, error)
	GetList(param entity.MidtransTransactionParam) ([]entity.MidtransTransaction, error)
	Get(param entity.MidtransTransactionParam) (entity.MidtransTransaction, error)
	Update(selectParam entity.MidtransTransactionParam, updateParam entity.UpdateMidtransTransactionParam) error
	Delete(param entity.MidtransTransactionParam) error
}

type transaction struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	t := &transaction{
		db: db,
	}

	return t
}

func (t *transaction) Create(transaction entity.MidtransTransaction) (entity.MidtransTransaction, error) {
	if err := t.db.Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (t *transaction) GetList(param entity.MidtransTransactionParam) ([]entity.MidtransTransaction, error) {
	transactions := []entity.MidtransTransaction{}

	if err := t.db.Where(param).Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (t *transaction) Get(param entity.MidtransTransactionParam) (entity.MidtransTransaction, error) {
	transaction := entity.MidtransTransaction{}

	if err := t.db.Where(param).First(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (t *transaction) Update(selectParam entity.MidtransTransactionParam, updateParam entity.UpdateMidtransTransactionParam) error {
	if err := t.db.Model(&selectParam).Updates(entity.MidtransTransaction{
		Status: updateParam.Status,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (t *transaction) Delete(param entity.MidtransTransactionParam) error {
	if err := t.db.Where(param).Delete(&entity.MidtransTransaction{}).Error; err != nil {
		return err
	}

	return nil
}
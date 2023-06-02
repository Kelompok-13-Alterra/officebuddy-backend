package transaction

import (
	"go-clean/src/business/entity"
	"time"

	"gorm.io/gorm"
)

type Interface interface {
	Create(transaction entity.Transaction) (entity.Transaction, error)
	GetList(param entity.TransactionParam) ([]entity.Transaction, error)
	GetListBooked(param entity.TransactionParam) ([]entity.Transaction, error)
	GetListHistoryBooked(param entity.TransactionParam) ([]entity.Transaction, error)
	Get(param entity.TransactionParam) (entity.Transaction, error)
	GetAvaibility(param entity.TransactionParam) (entity.Transaction, error)
	Update(selectParam entity.TransactionParam, updateParam entity.UpdateTransactionParam) error
	Delete(param entity.TransactionParam) error
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

func (t *transaction) Create(transaction entity.Transaction) (entity.Transaction, error) {
	if err := t.db.Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (t *transaction) GetList(param entity.TransactionParam) ([]entity.Transaction, error) {
	transactions := []entity.Transaction{}

	if err := t.db.Where(param).Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (t *transaction) GetListBooked(param entity.TransactionParam) ([]entity.Transaction, error) {
	transactions := []entity.Transaction{}

	if err := t.db.Where("start >= ? and user_id = ?", time.Now(), param.UserID).Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (t *transaction) GetListHistoryBooked(param entity.TransactionParam) ([]entity.Transaction, error) {
	transactions := []entity.Transaction{}

	if err := t.db.Where("start <= ? and user_id = ?", time.Now(), param.UserID).Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (t *transaction) GetAvaibility(param entity.TransactionParam) (entity.Transaction, error) {
	transaction := entity.Transaction{}

	if err := t.db.Where("((? BETWEEN start AND end) OR (? BETWEEN start AND end)) AND office_id = ?", param.Start, param.End, param.OfficeID).Limit(1).Find(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (t *transaction) Get(param entity.TransactionParam) (entity.Transaction, error) {
	transaction := entity.Transaction{}

	if err := t.db.Where(param).First(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (t *transaction) Update(selectParam entity.TransactionParam, updateParam entity.UpdateTransactionParam) error {
	if err := t.db.Model(&selectParam).Updates(entity.Transaction{
		Status: updateParam.Status,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (t *transaction) Delete(param entity.TransactionParam) error {
	if err := t.db.Where(param).Delete(&entity.Transaction{}).Error; err != nil {
		return err
	}

	return nil
}

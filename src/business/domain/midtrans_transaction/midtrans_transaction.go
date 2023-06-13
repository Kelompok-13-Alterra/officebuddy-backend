package midtrans_transaction

import (
	"fmt"
	"go-clean/src/business/entity"
	"time"

	"gorm.io/gorm"
)

type Interface interface {
	Create(transaction entity.MidtransTransaction) (entity.MidtransTransaction, error)
	GetList(param entity.MidtransTransactionParam) ([]entity.MidtransTransaction, error)
	Get(param entity.MidtransTransactionParam) (entity.MidtransTransaction, error)
	Update(selectParam entity.MidtransTransactionParam, updateParam entity.UpdateMidtransTransactionParam) error
	Delete(param entity.MidtransTransactionParam) error
	GetListWithPaginationByIDs(param entity.MidtransTransactionParam) ([]entity.MidtransTransaction, error)
	GetMidtransTransactionToday() ([]entity.MidtransTransaction, error)
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
	if err := t.db.Model(entity.MidtransTransaction{}).Where(selectParam).Updates(updateParam).Error; err != nil {
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

func (t *transaction) GetListWithPaginationByIDs(param entity.MidtransTransactionParam) ([]entity.MidtransTransaction, error) {
	transactions := []entity.MidtransTransaction{}

	if err := t.db.Limit(param.Limit).Offset(param.Offset).Where(param).Order(param.OrderBy).Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (t *transaction) GetMidtransTransactionToday() ([]entity.MidtransTransaction, error) {
	midtransTransaction := []entity.MidtransTransaction{}
	currentDate := time.Now().Format("2006-01-02")
	if err := t.db.Where("created_at LIKE ? AND status = ?", fmt.Sprintf("%%%s%%", currentDate), "success").Find(&midtransTransaction).Error; err != nil {
		return midtransTransaction, err
	}

	return midtransTransaction, nil
}

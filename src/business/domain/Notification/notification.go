package notification

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type interface interface {
	Create(notification entity.notification) (entity.notification, error)
	GetList(param entity.notificationParam) ([]entity.notification, error)
	Get(param entity.notificationParam) (entity.notification, error)
	Update(selectParam entity.notificationParam, updateParam entity.UpdatenotificationParam) error
	Delete(param entity.notificationParam) error
}

type notification struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	t := &notification{
		db: db,
	}

	return t
}

func (t *notification) Create(notification entity.notification) (entity.notification, error) {
	if err := t.db.Create(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (t *notification) GetList(param entity.notificationParam) ([]entity.notification, error) {
	notification := []entity.notification{}

	if err := t.db.Where(param).Find(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (t *notification) Get(param entity.notificationParam) (entity.notification, error) {
	notification := entity.notification{}

	if err := t.db.Where(param).First(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (t *notification) Update(selectParam entity.notificationParam, updateParam entity.UpdatenotificationParam) error {
	if err := t.db.Model(&selectParam).Updates(entity.notification{
		Status: updateParam.Status,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (t *notification) Delete(param entity.notificationParam) error {
	if err := t.db.Where(param).Delete(&entity.notification{}).Error; err != nil {
		return err
	}

	return nil
}

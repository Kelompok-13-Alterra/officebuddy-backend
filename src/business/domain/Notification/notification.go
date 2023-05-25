package notification

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(notification entity.Notification) (entity.Notification, error)
	GetList(param entity.NotificationParam) ([]entity.Notification, error)
	Get(param entity.NotificationParam) (entity.Notification, error)
	Delete(param entity.NotificationParam) error
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

func (t *notification) Create(Notification entity.Notification) (entity.Notification, error) {
	if err := t.db.Create(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (t *notification) GetList(param entity.NotificationParam) ([]entity.Notification, error) {
	notification := []entity.notification{}

	if err := t.db.Where(param).Find(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (t *notification) Get(param entity.NotificationParam) (entity.Notification, error) {
	notification := entity.notification{}

	if err := t.db.Where(param).First(&notification).Error; err != nil {
		return notification, err
	}

	return notification, nil
}

func (t *notification) Delete(param entity.NotificationParam) error {
	if err := t.db.Where(param).Delete(&entity.Notification{}).Error; err != nil {
		return err
	}

	return nil
}

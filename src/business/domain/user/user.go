package user

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(user entity.User) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	GetById(id uint) (entity.User, error)
	Update(selectParam entity.UpdateUserParam, updateParam entity.UpdateUserParam) error
}

type user struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	a := &user{
		db: db,
	}

	return a
}

func (a *user) Create(user entity.User) (entity.User, error) {
	if err := a.db.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (a *user) GetByEmail(email string) (entity.User, error) {
	user := entity.User{}

	if err := a.db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (a *user) GetById(id uint) (entity.User, error) {
	user := entity.User{}

	if err := a.db.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (o *user) Update(selectParam entity.UpdateUserParam, updateParam entity.UpdateUserParam) error {
	if err := o.db.Model(entity.User{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}

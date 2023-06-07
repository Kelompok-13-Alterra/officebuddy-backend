package user

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(user entity.User) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	GetById(id uint) (entity.User, error)
	GetList(param entity.UserParam) ([]entity.User, error)
	Update(selectParam entity.UpdateUserParam, updateParam entity.UpdateUserParam) error
	Delete(param entity.UserParam) error
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

func (r *user) GetList(param entity.UserParam) ([]entity.User, error) {
	users := []entity.User{}

	if err := r.db.Where(param).Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (o *user) Update(selectParam entity.UpdateUserParam, updateParam entity.UpdateUserParam) error {
	if err := o.db.Model(entity.User{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}

func (t *user) Delete(param entity.UserParam) error {
	if err := t.db.Where(param).Delete(&entity.UserParam{}).Error; err != nil {
		return err
	}

	return nil
}

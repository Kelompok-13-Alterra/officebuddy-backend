package rating

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(rating entity.Rating) (entity.Rating, error)
	GetList(param entity.RatingParam) ([]entity.Rating, error)
	Get(param entity.RatingParam) (entity.Rating, error)
	GetCount() (int64, error)
	GetCountInByID(ids []uint) (int64, error)
	Delete(param entity.RatingParam) error
}

type rating struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	r := &rating{
		db: db,
	}

	return r
}

func (r *rating) Create(rating entity.Rating) (entity.Rating, error) {
	if err := r.db.Create(&rating).Error; err != nil {
		return rating, err
	}

	return rating, nil
}

func (r *rating) GetList(param entity.RatingParam) ([]entity.Rating, error) {
	ratings := []entity.Rating{}

	if err := r.db.Where(param).Find(&ratings).Error; err != nil {
		return ratings, err
	}

	return ratings, nil
}

func (r *rating) Get(param entity.RatingParam) (entity.Rating, error) {
	rating := entity.Rating{}

	if err := r.db.Where(param).First(&rating).Error; err != nil {
		return rating, err
	}

	return rating, nil
}

func (r *rating) GetCountInByID(ids []uint) (int64, error) {
	var result int64
	if err := r.db.Model(entity.Rating{}).Where("office_id IN ?", ids).Count(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *rating) GetCount() (int64, error) {
	var result int64
	if err := r.db.Model(entity.Rating{}).Count(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (t *rating) Delete(param entity.RatingParam) error {
	if err := t.db.Where(param).Delete(&entity.Rating{}).Error; err != nil {
		return err
	}

	return nil
}

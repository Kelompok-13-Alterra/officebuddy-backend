package rating

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(rating entity.Rating) (entity.Rating, error)
	GetList(param entity.RatingParam) ([]entity.Rating, error)
	Get(param entity.RatingParam) (entity.Rating, error)
}

type rating struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	t := &rating{
		db: db,
	}

	return t
}

func (t *rating) Create(rating entity.Rating) (entity.Rating, error) {
	if err := t.db.Create(&rating).Error; err != nil {
		return rating, err
	}

	return rating, nil
}

func (t *rating) GetList(param entity.RatingParam) ([]entity.Rating, error) {
	ratings := []entity.Rating{}

	if err := t.db.Where(param).Find(&ratings).Error; err != nil {
		return ratings, err
	}

	return ratings, nil
}

func (t *rating) Get(param entity.RatingParam) (entity.Rating, error) {
	rating := entity.Rating{}

	if err := t.db.Where(param).First(&rating).Error; err != nil {
		return rating, err
	}

	return rating, nil
}

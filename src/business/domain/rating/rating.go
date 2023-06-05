package rating

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(rating entity.Rating) (entity.Rating, error)
	GetList(param entity.RatingParam) ([]entity.Rating, error)
	GetListCount(param entity.RatingParam) (int64, error)
	Get(param entity.RatingParam) (entity.Rating, error)
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

func (r *rating) GetListCount(param entity.RatingParam) (int64, error) {
	var count int64

	if err := r.db.Model(&entity.Rating{}).
		Where(&entity.Rating{
			TransactionID: param.TransactionID,
		}).
		Count(&count).Error; err != nil {
		return count, err
	}

	return count, nil
}

func (r *rating) Get(param entity.RatingParam) (entity.Rating, error) {
	rating := entity.Rating{}

	if err := r.db.Where(param).First(&rating).Error; err != nil {
		return rating, err
	}

	return rating, nil
}

func (t *rating) Delete(param entity.RatingParam) error {
	if err := t.db.Where(param).Delete(&entity.Rating{}).Error; err != nil {
		return err
	}

	return nil
}

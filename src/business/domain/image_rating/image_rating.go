package imagerating

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(imageRating entity.ImageRating) (entity.ImageRating, error)
	GetList(param entity.ImageRatingParam) ([]entity.ImageRating, error)
	Get(param entity.ImageRatingParam) (entity.ImageRating, error)
}

type imageRating struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	ir := &imageRating{
		db: db,
	}

	return ir
}

func (ir *imageRating) Create(imageRating entity.ImageRating) (entity.ImageRating, error) {
	if err := ir.db.Create(&imageRating).Error; err != nil {
		return imageRating, err
	}

	return imageRating, nil
}

func (ir *imageRating) GetList(param entity.ImageRatingParam) ([]entity.ImageRating, error) {
	imageRatings := []entity.ImageRating{}

	if err := ir.db.Where(param).Find(&imageRatings).Error; err != nil {
		return imageRatings, err
	}

	return imageRatings, nil
}

func (ir *imageRating) Get(param entity.ImageRatingParam) (entity.ImageRating, error) {
	imageRating := entity.ImageRating{}

	if err := ir.db.Where(param).First(&imageRating).Error; err != nil {
		return imageRating, err
	}

	return imageRating, nil
}

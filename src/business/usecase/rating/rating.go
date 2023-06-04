package rating

import (
	ratingDom "go-clean/src/business/domain/rating"
	"go-clean/src/business/entity"
	
)

type Interface interface {
	GetList(param entity.RatingParam)([]entity.Rating, error)
}

type rating struct {
	rating 	ratingDom.Interface
	
}

func Init(rt ratingDom.Interface) Interface {
	r := &rating{
		rating: rt,
	}

	return r
}

func (r *rating) GetList(param entity.RatingParam)([]entity.Rating, error) {
	var (
		ratings []entity.Rating
		err           error
	)

	ratings, err = r.rating.GetList(entity.RatingParam{})

	if err != nil {
		return ratings, err
	}

	return ratings, nil
}

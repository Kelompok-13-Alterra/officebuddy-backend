package rating

import (
	"context"
	ratingDom "go-clean/src/business/domain/rating"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateRatingParam) (entity.Rating, error)
}

type rating struct {
	rating ratingDom.Interface
	auth   auth.Interface
}

func Init(rd ratingDom.Interface, auth auth.Interface) Interface {
	o := &rating{
		rating: rd,
		auth:   auth,
	}

	return o
}

func (r *rating) Create(ctx context.Context, param entity.CreateRatingParam) (entity.Rating, error) {
	rating := entity.Rating{}

	user, err := r.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return rating, err
	}

	rating, err = r.rating.Create(entity.Rating{
		UserID:      user.User.ID,
		OfficeID:    param.OfficeID,
		Star:        param.Star,
		Tags:        param.Tags,
		Description: param.Description,
	})

	if err != nil {
		return rating, err
	}

	return rating, nil
}

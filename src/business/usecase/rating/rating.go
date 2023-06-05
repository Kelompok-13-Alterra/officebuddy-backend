package rating

import (
	"context"
	"encoding/json"
	"errors"
	ratingDom "go-clean/src/business/domain/rating"
	transactionDom "go-clean/src/business/domain/transaction"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
	"time"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateRatingParam) (entity.Rating, error)
	GetList(param entity.RatingParam) ([]entity.Rating, error)
}

type rating struct {
	rating      ratingDom.Interface
	transaction transactionDom.Interface
	auth        auth.Interface
}

func Init(rd ratingDom.Interface, td transactionDom.Interface, auth auth.Interface) Interface {
	r := &rating{
		transaction: td,
		rating:      rd,
		auth:        auth,
	}

	return r
}

func (r *rating) Create(ctx context.Context, param entity.CreateRatingParam) (entity.Rating, error) {
	rating := entity.Rating{}

	user, err := r.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return rating, err
	}

	transaction, err := r.transaction.Get(entity.TransactionParam{
		ID: param.TransactionID,
	})
	if err != nil {
		return rating, err
	}

	if time.Now().After(transaction.End) {
		return rating, errors.New("transaksi belum selesai")
	}

	tag, err := json.Marshal(param.Tags)
	if err != nil {
		return rating, err
	}

	rating, err = r.rating.Create(entity.Rating{
		UserID:        user.User.ID,
		OfficeID:      transaction.OfficeID,
		TransactionID: param.TransactionID,
		Star:          param.Star,
		Tags:          string(tag),
		Description:   param.Description,
	})

	if err != nil {
		return rating, err
	}

	return rating, nil
}

func (r *rating) GetList(param entity.RatingParam) ([]entity.Rating, error) {
	var (
		ratings []entity.Rating
		err     error
	)

	ratings, err = r.rating.GetList(entity.RatingParam{})

	if err != nil {
		return ratings, err
	}

	return ratings, nil
}

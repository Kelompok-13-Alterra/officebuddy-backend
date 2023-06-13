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
	GetList(param entity.RatingParam) ([]entity.RatingResponse, error)
	Get(param entity.RatingParam) (entity.RatingResponse, error)
	Delete(param entity.RatingParam) error
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

	rating, err = r.rating.Get(entity.RatingParam{
		TransactionID: param.TransactionID,
		OfficeID:      transaction.OfficeID,
	})
	if err == nil {
		return rating, errors.New("kamu sudah melakukan review")
	}

	if time.Now().Before(transaction.End) {
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

func (r *rating) GetList(param entity.RatingParam) ([]entity.RatingResponse, error) {
	var (
		ratings []entity.Rating
		result  []entity.RatingResponse
		err     error
	)

	ratings, err = r.rating.GetList(entity.RatingParam{})

	for _, r := range ratings {
		tags := []string{}
		_ = json.Unmarshal([]byte(r.Tags), &tags)
		result = append(result, entity.RatingResponse{
			ID:            r.ID,
			UserID:        r.UserID,
			OfficeID:      r.OfficeID,
			TransactionID: r.TransactionID,
			Star:          r.Star,
			Tags:          tags,
			Description:   r.Description,
			CreatedAt:     r.CreatedAt,
		})
	}

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *rating) Get(param entity.RatingParam) (entity.RatingResponse, error) {
	rating, err := r.rating.Get(param)
	if err != nil {
		return entity.RatingResponse{}, err
	}

	tags := []string{}
	_ = json.Unmarshal([]byte(rating.Tags), &tags)
	result := entity.RatingResponse{
		ID:            rating.ID,
		UserID:        rating.UserID,
		OfficeID:      rating.OfficeID,
		TransactionID: rating.TransactionID,
		Star:          rating.Star,
		Tags:          tags,
		Description:   rating.Description,
		CreatedAt:     rating.CreatedAt,
	}

	return result, nil
}

func (o *rating) Delete(param entity.RatingParam) error {
	if err := o.rating.Delete(entity.RatingParam{
		ID: param.ID,
	}); err != nil {
		return err
	}

	return nil
}

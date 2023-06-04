package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *rest) GetRatingList(ctx *gin.Context) {
	rating, err := r.uc.Rating.GetList(entity.RatingParam{})
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully get rating list", rating)
}

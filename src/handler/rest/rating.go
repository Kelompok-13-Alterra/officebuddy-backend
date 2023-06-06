package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create Rating
// @Description Create New Rating
// @Security BearerAuth
// @Tags Rating
// @Param rating_id path integer true "rating id"
// @Param rating body entity.CreateRatingParam true "rating info"
// @Produce json
// @Success 200 {object} entity.Response{data=entity.Rating{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/rating/{transaction_id} [POST]
func (r *rest) CreateRating(ctx *gin.Context) {
	var inputParam entity.CreateRatingParam
	if err := ctx.ShouldBindJSON(&inputParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := ctx.ShouldBindUri(&inputParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	rating, err := r.uc.Rating.Create(ctx.Request.Context(), inputParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully created rating", rating)
}

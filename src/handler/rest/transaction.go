package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Transaction Booked List
// @Description Get Transaction Booked List
// @Security BearerAuth
// @Tags Transaction
// @Produce json
// @Success 200 {object} entity.Response{data=entity.Transaction}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction/booked [GET]
func (r *rest) GetTransactionBookedList(ctx *gin.Context) {
	var param entity.TransactionParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := ctx.ShouldBindQuery(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.Transaction.GetListBooked(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get transactions booked list", result)
}

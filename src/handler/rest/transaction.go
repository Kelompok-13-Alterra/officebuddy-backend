package rest

import (
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
	result, err := r.uc.Transaction.GetListBooked(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get transactions booked list", result)
}

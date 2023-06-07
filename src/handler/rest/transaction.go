package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create Order
// @Description Create New Order
// @Security BearerAuth
// @Tags Transaction
// @Param office_id path integer true "office id"
// @Param transaction body entity.CreateTransactionParam true "transaction info"
// @Produce json
// @Success 200 {object} entity.Response{data=int}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction/office/{office_id}/book [POST]
func (r *rest) CreateOrder(ctx *gin.Context) {
	var inputParam entity.CreateTransactionParam
	if err := ctx.ShouldBindJSON(&inputParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := ctx.ShouldBindUri(&inputParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	id, err := r.uc.Transaction.Create(ctx.Request.Context(), inputParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully created new order", gin.H{"id_transaction": id})
}

// @Summary Reschedule Transaction
// @Description Reschedule transaction
// @Security BearerAuth
// @Tags Transaction
// @Param transaction_id path integer true "tranasction id"
// @Param transaction body entity.InputUpdateTransactionParam true "transaction info"
// @Produce json
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction/{transaction_id}/reschedule [PUT]
func (r *rest) RescheduleBooked(ctx *gin.Context) {
	var inputParam entity.InputUpdateTransactionParam
	if err := ctx.ShouldBindJSON(&inputParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	var selectParam entity.TransactionParam
	if err := ctx.ShouldBindUri(&selectParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.Transaction.RescheduleBooked(ctx.Request.Context(), inputParam, selectParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully reschedule transaction", nil)
}

// @Summary Get Transaction Booked List
// @Description Get Transaction Booked List
// @Security BearerAuth
// @Tags Transaction
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Transaction}
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

// @Summary Get Transaction History Booked List
// @Description Get Transaction History Booked List
// @Security BearerAuth
// @Tags Transaction
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Transaction}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction/history [GET]
func (r *rest) GetTransactionHistoryBookedList(ctx *gin.Context) {
	result, err := r.uc.Transaction.GetListHistoryBooked(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get transactions history booked list", result)
}

// @Summary Get Transaction List
// @Description Get Transaction List for Admin
// @Security BearerAuth
// @Tags Transaction
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Transaction}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/transaction [GET]
func (r *rest) GetTransactionList(ctx *gin.Context) {
	transaction, err := r.uc.Transaction.GetTransactionList(entity.TransactionParam{})
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully get transaction list", transaction)
}

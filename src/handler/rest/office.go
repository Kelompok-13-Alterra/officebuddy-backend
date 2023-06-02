package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Summary Create Office
// @Description Create new Office
// @Security BearerAuth
// @Tags Office
// @Produce json
// @Param office body entity.CreateOfficeParam true "office info"
// @Param type query string true "type" Enums(office, coworking)
// @Success 200 {object} entity.Response{data=entity.Office{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/office [POST]
func (r *rest) CreateOffice(ctx *gin.Context) {
	var inputParam entity.CreateOfficeParam
	if err := ctx.ShouldBindJSON(&inputParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	var typeParam entity.OfficeTypeParam
	if err := ctx.ShouldBindWith(&typeParam, binding.Query); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	inputParam.Type = typeParam.Type

	office, err := r.uc.Office.Create(inputParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully create new office", office)
}

// @Summary Get Office List
// @Description Get Office List
// @Security BearerAuth
// @Tags Office
// @Produce json
// @Param name query string false "name"
// @Success 200 {object} entity.Response{data=[]entity.Office{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/office [GET]
func (r *rest) GetOfficeList(ctx *gin.Context) {
	var officeParam entity.OfficeParam
	if err := ctx.ShouldBindWith(&officeParam, binding.Query); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	office, err := r.uc.Office.GetList(officeParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get office list", office)
}

// @Summary Get Office Detail
// @Description Get Office Detail
// @Security BearerAuth
// @Tags Office
// @Produce json
// @Param office_id path integer true "office id"
// @Success 200 {object} entity.Response{data=entity.Office{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/office/{office_id} [GET]
func (r *rest) GetOffice(ctx *gin.Context) {
	var param entity.OfficeParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	office, err := r.uc.Office.Get(param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfull get office detail", office)
}

// @Summary Update Office
// @Description Update a Office
// @Security BearerAuth
// @Tags Office
// @Produce json
// @Param office body entity.UpdateOfficeParam true "office info"
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/office/{office_id} [PUT]
func (r *rest) UpdateOffice(ctx *gin.Context) {
	var updateParam entity.UpdateOfficeParam
	if err := ctx.ShouldBindJSON(&updateParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	var selectParam entity.OfficeParam
	if err := ctx.ShouldBindUri(&selectParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.Office.Update(selectParam, updateParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully update office", nil)
}

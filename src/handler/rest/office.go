package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

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

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully get office list", office)
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

package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (r *rest) GetOfficeList(ctx *gin.Context) {
	var officeParam entity.OfficeParam
	if err := ctx.ShouldBindWith(&officeParam, binding.Query); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	offices, err := r.uc.Office.GetList(officeParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get office list", offices)
}

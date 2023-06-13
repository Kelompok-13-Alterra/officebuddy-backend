package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Summary Get Dashboard Widget
// @Description Get Dashboard Widget for Admin
// @Security BearerAuth
// @Tags Widget Analytic
// @Produce json
// @Success 200 {object} entity.Response{data=entity.DashboardWidgetResult}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/admin/dashboard-widget [GET]
func (r *rest) GetDashboardWidget(ctx *gin.Context) {
	result, err := r.uc.WidgetDashboard.GetDashboardWidget(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get dashboard widget", result)
}

// @Summary Get Office Widget
// @Description Get Office or Coworking Widget for Admin
// @Security BearerAuth
// @Tags Widget Analytic
// @Produce json
// @Param type query string true "type" Enums(office, coworking)
// @Success 200 {object} entity.Response{data=entity.DashboardWidgetResult}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/admin/office-widget [GET]
func (r *rest) GetOfficeWidget(ctx *gin.Context) {
	var param entity.OfficeWidgetParam
	if err := ctx.ShouldBindWith(&param, binding.Query); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := r.uc.WidgetDashboard.GetOfficeWidget(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get office widget", result)
}

// @Summary Get Revenue Widget
// @Description Get Revenue Widget for Admin
// @Security BearerAuth
// @Tags Widget Analytic
// @Produce json
// @Success 200 {object} entity.Response{data=entity.RevenueWidgetResult}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/admin/revenue-widget [GET]
func (r *rest) GetRevenueWidget(ctx *gin.Context) {
	result, err := r.uc.WidgetDashboard.GetRevenueWidget(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get revenue widget", result)
}

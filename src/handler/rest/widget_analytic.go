package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

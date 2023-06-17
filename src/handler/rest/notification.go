package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Notification List
// @Description Get Notification List by User Logged In
// @Security BearerAuth
// @Tags Notification
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Notification}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/notification [GET]
func (r *rest) GetNotificationList(ctx *gin.Context) {
	notification, err := r.uc.Notification.GetList(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully get notification list", notification)
}

func (r *rest) MarkAsRead(ctx *gin.Context) {
	notification, err := r.uc.Notification.MarkAsRead(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully Mark As Read list", notification)
}

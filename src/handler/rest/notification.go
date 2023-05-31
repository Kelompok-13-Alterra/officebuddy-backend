package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *rest) GetNotificationList(ctx *gin.Context) {
	notification, err := r.uc.Notification.GetList(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully get notification list", notification)
}

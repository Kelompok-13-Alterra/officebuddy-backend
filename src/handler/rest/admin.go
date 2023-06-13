package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Rating List
// @Description Get Rating List
// @Security BearerAuth
// @Tags Rating
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Rating{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/rating [GET]
func (r *rest) GetRatingList(ctx *gin.Context) {
	rating, err := r.uc.Rating.GetList(entity.RatingParam{})
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully get rating list", rating)
}

// @Summary Get Rating Detail
// @Description Get Rating Detail
// @Security BearerAuth
// @Tags Rating
// @Produce json
// @Param rating_id path integer true "rating id"
// @Success 200 {object} entity.Response{data=entity.Rating{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/rating/{rating_id} [GET]
func (r *rest) GetRating(ctx *gin.Context) {
	var param entity.RatingParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	rating, err := r.uc.Rating.Get(param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfull get rating detail", rating)
}

// @Summary Get User List
// @Description Get User List
// @Security BearerAuth
// @Tags User
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.User{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/user [GET]
func (r *rest) GetUserList(ctx *gin.Context) {
	user, err := r.uc.User.GetUserList(entity.UserParam{})
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "successfully get user list", user)
}

// @Summary Update User
// @Description Update a User
// @Security BearerAuth
// @Tags User
// @Produce json
// @Param user_id path integer true "user id"
// @Param user body entity.UpdateUserParam true "user info"
// @Success 200 {object} entity.Response{}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/user/edit/{user_id} [PUT]
func (r *rest) UpdateUserByAdmin(ctx *gin.Context) {
	var param entity.UserParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	var updateParam entity.UpdateUserParam
	if err := ctx.ShouldBindJSON(&updateParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, err)
		return
	}

	err := r.uc.User.UpdateByAdmin(param, updateParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully update user", nil)
}

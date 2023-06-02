package entity

import (
	"go-clean/src/lib/auth"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string `json:"-"`
	Role     int
	IsVerify bool
}

type CreateUserParam struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
	Name     string `binding:"required"`
}

type LoginUserParam struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type UpdateUserParam struct {
	Name     string
	Email    string
	Password string
}

type UserParam struct {
	ID uint
}

func (u *User) ConvertToAuthUser() auth.User {
	return auth.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Role:     u.Role,
		IsVerify: u.IsVerify,
	}
}

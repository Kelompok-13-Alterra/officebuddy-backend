package user

import (
	"context"
	"errors"
	userDom "go-clean/src/business/domain/user"
	"go-clean/src/business/entity"
	auth "go-clean/src/lib/auth"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Interface interface {
	Create(params entity.CreateUserParam) (entity.User, error)
	Login(params entity.LoginUserParam) (string, error)
	LoginAdmin(params entity.LoginUserParam) (string, error)
	Get(ctx context.Context) (entity.User, error)
	GetById(id uint) (entity.User, error)
	Update(ctx context.Context, inputParam entity.UpdateUserParam) error
}

type user struct {
	user userDom.Interface
	auth auth.Interface
}

func Init(ad userDom.Interface, auth auth.Interface) Interface {
	a := &user{
		user: ad,
		auth: auth,
	}

	return a
}

func (a *user) Create(params entity.CreateUserParam) (entity.User, error) {
	user := entity.User{
		Email:    params.Email,
		Name:     params.Name,
		Role:     auth.RoleUser,
		IsVerify: false,
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.Password = string(hashPass)

	newUser, err := a.user.Create(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (u *user) Get(ctx context.Context) (entity.User, error) {
	user, err := u.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.User{}, err
	}
	userProfile, err := u.user.GetById(user.User.ID)
	if err != nil {
		return userProfile, err
	}
	return userProfile, nil
}

func (a *user) GetById(id uint) (entity.User, error) {
	user, err := a.user.GetById(id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (a *user) Login(params entity.LoginUserParam) (string, error) {
	user, err := a.user.GetByEmail(params.Email)
	if err != nil {
		return "", err
	}

	if user.ID == 0 {
		return "", errors.New("user tidak ditemukan atau password tidak sesuai")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return "", errors.New("user tidak ditemukan atau password tidak sesuai")
	}

	token, err := a.auth.GenerateToken(user.ConvertToAuthUser())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *user) LoginAdmin(params entity.LoginUserParam) (string, error) {
	user, err := a.user.GetByEmail(params.Email)
	if err != nil {
		return "", err
	}

	if user.ID == 0 {
		return "", errors.New("user tidak ditemukan atau password tidak sesuai")
	}

	if user.Role != auth.RoleAdmin {
		return "", errors.New("user tidak ditemukan atau password tidak sesuai")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return "", errors.New("user tidak ditemukan atau password tidak sesuai")
	}

	token, err := a.auth.GenerateToken(user.ConvertToAuthUser())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *user) Update(ctx context.Context, inputParam entity.UpdateUserParam) error {
	user, err := u.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	if inputParam.DateBirthInput != "" {
		formatedDate, err := u.formatDate(inputParam.DateBirthInput)
		if err != nil {
			return err
		}
		inputParam.DateBirth = formatedDate
	}

	if err := u.user.Update(entity.UpdateUserParam{
		Email: user.User.Email,
	}, inputParam); err != nil {
		return err
	}

	return nil
}

func (u *user) formatDate(date string) (time.Time, error) {
	var formatedDate time.Time

	layoutFormat := "2006-01-02"
	formatedDate, err := time.Parse(layoutFormat, date)
	if err != nil {
		return formatedDate, err
	}

	return formatedDate, nil
}

package auth

import (
	"context"
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	userAuthInfo contextKey = "UserAuthInfo"
	RoleUser                = 0
	RoleAdmin               = 1
)

type Interface interface {
	SetUserAuthInfo(ctx context.Context, user User, token string) context.Context
	GetUserAuthInfo(ctx context.Context) (UserAuthInfo, error)
	GenerateToken(user User) (string, error)
}

type auth struct {
}

func Init() Interface {
	return &auth{}
}

func (a *auth) SetUserAuthInfo(ctx context.Context, user User, token string) context.Context {
	userAuth := UserAuthInfo{
		User:  user,
		Token: token,
	}

	return context.WithValue(ctx, userAuthInfo, userAuth)
}

func (a *auth) GetUserAuthInfo(ctx context.Context) (UserAuthInfo, error) {
	userContext := ctx.Value(userAuthInfo)
	user, ok := userContext.(UserAuthInfo)
	if !ok {
		return user, errors.New("failed to get user auth")
	}

	return user, nil
}

func (a *auth) GenerateToken(user User) (string, error) {
	claim := jwt.MapClaims{}
	claim["id"] = user.ID
	claim["email"] = user.Email
	claim["is_verify"] = user.IsVerify
	claim["role"] = user.Role

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

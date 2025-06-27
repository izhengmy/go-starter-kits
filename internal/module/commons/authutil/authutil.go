package authutil

import (
	"app/pkg/auth"

	"github.com/gin-gonic/gin"
)

type AuthenticationUser[T uint] struct {
	ID       T      `json:"id"`
	Username string `json:"username"`
}

var _ auth.Authenticatable[uint] = (*AuthenticationUser[uint])(nil)

func (u AuthenticationUser[T]) JWTIdentifierKey() string {
	return "id"
}

func (u AuthenticationUser[T]) JWTIdentifier() T {
	return u.ID
}

func (u AuthenticationUser[T]) JWTCustomClaims() auth.JWTCustomClaims {
	return auth.JWTCustomClaims{
		"username": u.Username,
	}
}

func GetAuthenticationUser(ctx *gin.Context) *AuthenticationUser[uint] {
	value, ok := ctx.Get("user")
	if !ok {
		return nil
	}

	user, ok := value.(*AuthenticationUser[uint])
	if !ok {
		return nil
	}

	return user
}

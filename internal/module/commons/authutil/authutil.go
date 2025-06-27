package authutil

import (
	"app/internal/errorx"
	pkgHTTP "app/internal/pkg/http"
	"app/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

const CtxKeyUser = "user"

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
	value, ok := ctx.Get(CtxKeyUser)
	if !ok {
		return nil
	}

	user, ok := value.(*AuthenticationUser[uint])
	if !ok {
		return nil
	}

	return user
}

func Unauthenticated(json *pkgHTTP.JSON) auth.UnauthenticatedFunc {
	return func(ctx *gin.Context, err error) {
		json.Fail(ctx, errorx.NewServiceError("登录已过期").WithCode(http.StatusUnauthorized))
		ctx.AbortWithStatus(http.StatusOK)
	}
}

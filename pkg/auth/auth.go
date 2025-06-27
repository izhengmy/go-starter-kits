package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWTCustomClaims map[string]any

type Authenticatable[T uint | int | string] interface {
	JWTIdentifierKey() string
	JWTIdentifier() T
	JWTCustomClaims() JWTCustomClaims
}

type RetrieveTokenFunc func(ctx *gin.Context) string

type UnauthenticatedFunc func(ctx *gin.Context, err error)

var defaultRetrieveTokenFunc RetrieveTokenFunc = func(ctx *gin.Context) string {
	return ctx.GetHeader("Authorization")
}

var defaultUnauthenticatedFunc UnauthenticatedFunc = func(ctx *gin.Context, err error) {
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

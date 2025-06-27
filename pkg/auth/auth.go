package auth

import (
	"app/pkg/jwtx"
	"errors"
	"net/http"
	"strings"

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

func GinJWTAuthHandlerFunc[I uint | int | string, T Authenticatable[I]](auth *T, opts ...Option) gin.HandlerFunc {
	opts = append(
		opts,
		WithRetrieveTokenFunc(func(ctx *gin.Context) string {
			prefix := jwtx.TokenType + " "
			token := ctx.GetHeader("Authorization")
			if strings.HasPrefix(token, prefix) {
				return strings.TrimPrefix(token, prefix)
			} else {
				return ""
			}
		}),
	)
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}

	return func(ctx *gin.Context) {
		skip := false
		authed := false

		if len(o.skips) > 0 {
			skips := o.skips
			for _, path := range skips {
				if path == ctx.Request.URL.Path {
					skip = true
					break
				}
			}
		}

		retrieveTokenFunc := defaultRetrieveTokenFunc
		if o.retrieveTokenFunc != nil {
			retrieveTokenFunc = o.retrieveTokenFunc
		}

		unauthenticatedFunc := defaultUnauthenticatedFunc
		if o.unauthenticatedFunc != nil {
			unauthenticatedFunc = o.unauthenticatedFunc
		}

		token := retrieveTokenFunc(ctx)
		if token == "" {
			if !skip {
				unauthenticatedFunc(ctx, errors.New("token is empty"))
				return
			}
		}

		err := ParseJWTToken(token, auth)
		if err != nil {
			if !skip {
				unauthenticatedFunc(ctx, err)
				return
			}
		} else {
			authed = true
		}

		if authed {
			ctx.Set(o.ctxKey, auth)
		} else {
			ctx.Set(o.ctxKey, nil)
		}

		ctx.Next()
	}
}

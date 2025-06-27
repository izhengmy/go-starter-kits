package auth

import (
	"app/pkg/jwtx"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var reservedClaims = map[string]bool{
	"iss": true,
	"sub": true,
	"aud": true,
	"exp": true,
	"nbf": true,
	"iat": true,
	"jti": true,
}

func GenerateJWTToken[I uint | int | string, T Authenticatable[I]](auth T) (string, error) {
	claims := jwt.MapClaims{
		"sub": auth.JWTIdentifier(),
	}

	for k, v := range auth.JWTCustomClaims() {
		if !reservedClaims[k] {
			claims[k] = v
		}
	}

	return jwtx.GenerateToken(claims)
}

func ParseJWTToken[I uint | int | string, T Authenticatable[I]](token string, auth *T) error {
	claims, err := jwtx.ParseToken(token)
	if err != nil {
		return err
	}

	customClaims := JWTCustomClaims{}
	customClaims[(*auth).JWTIdentifierKey()] = claims["sub"]
	for k, v := range claims {
		if !reservedClaims[k] {
			customClaims[k] = v
		}
	}

	data, err := json.Marshal(customClaims)
	if err != nil {
		return fmt.Errorf("failed to marshal claims: %w", err)
	}

	if err := json.Unmarshal(data, auth); err != nil {
		return fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return nil
}

func GinJWTAuthHandlerFunc[I uint | int | string, T Authenticatable[I]](auth *T, abort bool, opts ...Option) gin.HandlerFunc {
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
		authed := false

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
			if abort {
				unauthenticatedFunc(ctx, errors.New("token is empty"))
				return
			}
		}

		err := ParseJWTToken(token, auth)
		if err != nil {
			if abort {
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

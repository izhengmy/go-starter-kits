package auth

import (
	"app/pkg/jwtx"
	"encoding/json"
	"fmt"

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

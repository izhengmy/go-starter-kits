package jwtx

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const TokenType = "Bearer"

func GenerateToken(claims jwt.MapClaims) (string, error) {
	issuer := viper.GetString("jwt.issuer")
	ttl := viper.GetInt("jwt.ttl")
	algo := viper.GetString("jwt.algo")
	secret := byteSecret()

	now := time.Now()
	expiresAt := now.Add(time.Duration(ttl) * time.Minute)

	method := jwt.SigningMethodHS256
	switch algo {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS512":
		method = jwt.SigningMethodHS512
	}

	claims["iss"] = issuer
	claims["exp"] = jwt.NewNumericDate(expiresAt)
	claims["iat"] = jwt.NewNumericDate(now)

	token := jwt.NewWithClaims(method, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ParseToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return byteSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func byteSecret() []byte {
	return []byte(viper.GetString("jwt.secret"))
}

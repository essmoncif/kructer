package core

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
	ID    string `json:"id"`
}

func NewJWTClaims(login, id string) *JWTClaims {
	return &JWTClaims{
		Login: login,
		ID:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
}

func (jwtc *JWTClaims) GenerateToken(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtc)
	return token.SignedString([]byte(secret))
}

package config

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type jwtService struct {
	secretKey string
	issuer    string
}

type JWTService interface {
	ValidateToken(ctx context.Context, tokenString string) (token *jwt.Token, err error)
}

func NewJWTService(secretKey string, apiUrl string) JWTService {
	return &jwtService{
		secretKey: secretKey,
		issuer:    apiUrl,
	}
}

func (s *jwtService) ValidateToken(ctx context.Context, tokenString string) (token *jwt.Token, err error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}

package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/peang/bukabengkel-api-go/src/models"
)

type TokenPayload struct {
	jwt.StandardClaims
	Payload Payload
	Scope   int `json:"scope"`
}

type Payload struct {
	ID            string
	FirstName     string
	LastName      interface{}
	Email         models.Email
	Mobile        models.Mobile
	Status        float64
	StoreRole     float64
	StoreID       string
	StoreName     string
	StoreType     float64
	StoreTypeName string
	StoreLocation models.Location
}

type jwtService struct {
	secretKey string
	issuer    string
}

type JWTService interface {
	ValidateToken(ctx context.Context, tokenString string) (token *jwt.Token, err error)
	GetTokenInfo(ctx context.Context, tokenString string) (tokenInfo TokenPayload, err error)
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

func (s *jwtService) GetTokenInfo(ctx context.Context, tokenString string) (tokenInfo TokenPayload, err error) {
	payload := &TokenPayload{}
	token, err := jwt.ParseWithClaims(tokenString, payload, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return TokenPayload{}, err
	}

	claim := token.Claims.(*TokenPayload)
	if claim.ExpiresAt < time.Now().Unix() {
		return TokenPayload{}, errors.New("token expired")
	}

	return *payload, nil
}

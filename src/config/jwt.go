package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/peang/bukabengkel-api-go/src/domain/entity"
	vo "github.com/peang/bukabengkel-api-go/src/domain/value_objects"
)

type TokenPayload struct {
	Issuer   string  `json:"iss"`
	Subject  string  `json:"sub"`
	Audience string  `json:"aud"`
	Nbf      int64   `json:"nbf"`
	Iat      int64   `json:"iat"`
	Exp      int64   `json:"exp"`
	Payload  Payload `json:"payload"`
	Scope    string  `json:"scope"`
}

type Payload struct {
	ID            string                `json:"id"`
	FirstName     string                `json:"firstname"`
	LastName      interface{}           `json:"lastname"`
	Email         vo.Email              `json:"email"`
	Mobile        vo.Mobile             `json:"mobile"`
	Status        float64               `json:"status"`
	StoreRole     float64               `json:"storeRole"`
	StoreID       string                `json:"storeId"`
	StoreName     string                `json:"storeName"`
	StoreType     float64               `json:"storeType"`
	StoreTypeName string                `json:"storeTypeName"`
	StoreLocation entity.LocationEntity `json:"storeLocation"`
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
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

	if err != nil {
		return TokenPayload{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		aud, _ := claims["aud"].(string)
		iat, _ := claims["iat"].(int64)
		exp, _ := claims["exp"].(int64)
		iss, _ := claims["iss"].(string)
		nbf, _ := claims["nbf"].(int64)
		sub, _ := claims["sub"].(string)
		payload, _ := claims["payload"].(map[string]interface{})
		scope, _ := claims["scope"].(string)

		// email := payload["email"]
		// _, cerr := strconv.Atoi(sub)
		// if cerr != nil {
		// 	err = errors.New("invalid token claims")
		// 	return
		// }

		var pyld = Payload{
			// 	// Email:     email.(vo.Email),
			// 	// Mobile:        payload["mobile"],
			Status:        payload["status"].(float64),
			StoreRole:     payload["storeRole"].(float64),
			StoreID:       payload["storeId"].(string),
			StoreName:     payload["storeName"].(string),
			StoreType:     payload["storeType"].(float64),
			StoreTypeName: payload["storeTypeName"].(string),
			// 	// StoreLocation: payload["storeLocation"].(string),
		}
		pyld.FirstName = payload["firstname"].(string)
		if value, ok := payload["lastname"].(string); ok {
			pyld.LastName = value
		}

		tokenInfo = TokenPayload{
			Issuer:   iss,
			Subject:  sub,
			Audience: aud,
			Nbf:      nbf,
			Iat:      iat,
			Exp:      exp,
			Payload:  pyld,
			Scope:    scope,
		}
	} else {
		err = errors.New("invalid token claims")
	}

	return
}

package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/peang/bukabengkel-api-go/src/config"
	utils "github.com/peang/bukabengkel-api-go/src/utils"
)

type Middleware struct {
	enforcer   *casbin.Enforcer
	logger     utils.Logger
	jwtservice config.JWTService
}

func NewMiddleware(enfocer *casbin.Enforcer, logger utils.Logger, jwtservice config.JWTService) *Middleware {
	return &Middleware{
		enforcer:   enfocer,
		logger:     logger,
		jwtservice: jwtservice,
	}
}

package driver

import (
	"github.com/urfave/negroni"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/sirupsen/logrus"
)

type MiddlewareRegistry interface {
	Common() *negroni.Negroni
	Logger() negroni.Handler
	Authorization() negroni.HandlerFunc
}

type DefaultMiddleware struct {
	common        *negroni.Negroni
	authorization negroni.HandlerFunc
	logger        *negronilogrus.Middleware

	r Registry
}

func NewDefaultMiddleware(r Registry) *DefaultMiddleware {
	return &DefaultMiddleware{
		r: r,
	}
}

func (m *DefaultMiddleware) Common() *negroni.Negroni {
	if m.common == nil {
		m.common = negroni.New(m.Logger())
	}
	return m.common
}

func (m *DefaultMiddleware) Logger() negroni.Handler {
	if m.logger == nil {
		logger := negronilogrus.NewMiddlewareFromLogger(m.r.Logger().(*logrus.Logger), "go-basic-structure")
		m.logger = logger
	}
	return m.logger
}

func (m *DefaultMiddleware) Authorization() negroni.HandlerFunc {
	if m.authorization == nil {
		authorization := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return m.r.AuthHandler().Authorize(token.Claims)
			},
			SigningMethod: jwt.SigningMethodRS256,
		})
		m.authorization = authorization.HandlerWithNext
	}

	return m.authorization
}

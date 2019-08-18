package driver

import (
	"github.com/gorilla/mux"
	"github.com/ory/herodot"
	"github.com/sawadashota/di-sample/auth"
	"github.com/sawadashota/di-sample/health"
	"github.com/sawadashota/di-sample/internal/admin"
	"github.com/sawadashota/di-sample/notification"
	"github.com/sirupsen/logrus"
)

type Registry interface {
	Writer() herodot.Writer
	Logger() logrus.FieldLogger

	Middleware() MiddlewareRegistry
	RegisterRoutes(router *mux.Router)

	admin.Registry

	HealthHandler() *health.Handler
	auth.Registry
	AuthHandler() *auth.Handler
	notification.Registry
	NotificationHandler() *notification.Handler
}

func NewDefaultRegistry() Registry {
	return NewRegistryMemory()
}

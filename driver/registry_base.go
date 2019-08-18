package driver

import (
	"github.com/gorilla/mux"
	"github.com/ory/herodot"
	"github.com/sawadashota/di-sample/auth"
	"github.com/sawadashota/di-sample/health"
	"github.com/sawadashota/di-sample/internal/admin"
	"github.com/sawadashota/di-sample/notification"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type RegistryBase struct {
	l      logrus.FieldLogger
	writer herodot.Writer

	adr admin.Repository

	heh *health.Handler
	aur auth.Repository
	auh *auth.Handler
	nor notification.Repository
	noh *notification.Handler

	m MiddlewareRegistry
	r Registry
}

func (rb *RegistryBase) Middleware() MiddlewareRegistry {
	if rb.m == nil {
		m := NewDefaultMiddleware(rb.r)
		rb.m = m
	}

	return rb.m
}

func (rb *RegistryBase) Writer() herodot.Writer {
	if rb.writer == nil {
		writer := herodot.NewJSONWriter(rb.Logger())
		rb.writer = writer
	}
	return rb.writer
}

func (rb *RegistryBase) Logger() logrus.FieldLogger {
	if rb.l == nil {
		l := logrus.New()
		l.SetFormatter(&logrus.JSONFormatter{})
		l.SetLevel(logrus.DebugLevel)
		rb.l = l
	}
	return rb.l
}

func (rb *RegistryBase) RegisterRoutes(router *mux.Router) {
	router.
		PathPrefix(health.PathPrefix).
		Handler(
			rb.Middleware().Common().With(
				negroni.Wrap(rb.HealthHandler().Router()),
			),
		)

	router.
		PathPrefix(auth.PathPrefix).
		Handler(
			rb.Middleware().Common().With(
				negroni.Wrap(rb.AuthHandler().Router()),
			),
		)

	router.
		PathPrefix(notification.PathPrefix).
		Handler(
			rb.Middleware().Common().With(
				rb.Middleware().Authorization(),
				negroni.Wrap(rb.NotificationHandler().Router()),
			),
		)
}

func (rb *RegistryBase) HealthHandler() *health.Handler {
	if rb.heh == nil {
		heh := health.NewHandler(rb.r)
		rb.heh = heh
	}
	return rb.heh
}

func (rb *RegistryBase) AuthHandler() *auth.Handler {
	if rb.auh == nil {
		auh := auth.NewHandler(rb.r)
		rb.auh = auh
	}
	return rb.auh
}

func (rb *RegistryBase) NotificationHandler() *notification.Handler {
	if rb.noh == nil {
		noh := notification.NewHandler(rb.r)
		rb.noh = noh
	}
	return rb.noh
}

func (rb *RegistryBase) with(r Registry) *RegistryBase {
	rb.r = r
	return rb
}

package health

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ory/herodot"
	"github.com/sirupsen/logrus"
)

const (
	PathPrefix = "/health"
	AlivePath  = "/health/alive"
)

type Handler struct {
	r InternalRegistry
}

type InternalRegistry interface {
	Writer() herodot.Writer
	Logger() logrus.FieldLogger
}

func NewHandler(r InternalRegistry) *Handler {
	return &Handler{
		r: r,
	}
}

func (h *Handler) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc(AlivePath, h.Alive).Methods(http.MethodGet)

	return router
}

func (h *Handler) Alive(w http.ResponseWriter, r *http.Request) {
	resp := &AliveResponse{Status: "ok"}
	h.r.Writer().Write(w, r, resp)
}

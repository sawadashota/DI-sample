package notification

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ory/herodot"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	PathPrefix             = "/notifications"
	ListNotificationPath   = "/notifications"
	GetNotificationPath    = "/notifications/{id}"
	AddNotificationPath    = "/notifications"
	UpdateNotificationPath = "/notifications/{id}"
	DeleteNotificationPath = "/notifications/{id}"
)

type Handler struct {
	r InternalRegistry
}

type InternalRegistry interface {
	Writer() herodot.Writer
	Logger() logrus.FieldLogger
	Registry
}

type Registry interface {
	NotificationRepository() Repository
}

func NewHandler(r InternalRegistry) *Handler {
	return &Handler{
		r: r,
	}
}

func (h *Handler) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc(ListNotificationPath, h.ListNotifications).Methods(http.MethodGet)
	router.HandleFunc(GetNotificationPath, h.GetNotification).Methods(http.MethodGet)
	router.HandleFunc(AddNotificationPath, h.AddNotification).Methods(http.MethodPost)
	router.HandleFunc(UpdateNotificationPath, h.UpdateNotification).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc(DeleteNotificationPath, h.DeleteNotification).Methods(http.MethodDelete, http.MethodOptions)

	return router
}

func (h *Handler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")

	const limit = 20
	offset := 0
	var err error
	if len(offsetStr) > 0 {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			h.r.Writer().WriteError(w, r, err)
			return
		}
	}

	ns, err := h.r.NotificationRepository().List(r.Context(), offset, limit)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	resp := &ListNotificationResponse{
		Notifications: ns,
	}

	h.r.Writer().Write(w, r, &resp)
}

func (h *Handler) GetNotification(w http.ResponseWriter, r *http.Request) {
	ps := mux.Vars(r)

	id, ok := ps["id"]
	if !ok {
		h.r.Writer().WriteError(w, r, errors.New("param id is not found"))
		return
	}

	n, err := h.r.NotificationRepository().Get(r.Context(), id)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, n)
}

func (h *Handler) AddNotification(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var p AddNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	n, err := NewNotification(uuid.New().String(), p.Title, p.Body, p.IsDraft, p.PublishAt)
	if err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.r.NotificationRepository().Add(r.Context(), n); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().WriteCreated(w, r, ListNotificationPath+"/"+n.ID, n)
}

func (h *Handler) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var p UpdateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	n, err := NewNotification(p.ID, p.Title, p.Body, p.IsDraft, p.PublishAt)
	if err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.r.NotificationRepository().Update(r.Context(), n); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	ps := mux.Vars(r)

	id, ok := ps["id"]
	if !ok {
		h.r.Writer().WriteError(w, r, errors.New("param id is not found"))
		return
	}

	err := h.r.NotificationRepository().Delete(r.Context(), id)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

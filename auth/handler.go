package auth

import (
	"context"
	"crypto"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/ory/herodot"
	"github.com/pkg/errors"
	"github.com/sawadashota/di-sample/internal/admin"
	"github.com/sirupsen/logrus"
)

const (
	PathPrefix = "/auth"
	LoginPath  = "/auth/login"
)

type Handler struct {
	r InternalRegistry
}

type InternalRegistry interface {
	Writer() herodot.Writer
	Logger() logrus.FieldLogger
	admin.Registry
	Registry
}

type Registry interface {
	JSONWebKeySetRepository() Repository
}

func NewHandler(r InternalRegistry) *Handler {
	return &Handler{
		r: r,
	}
}

func (h *Handler) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc(LoginPath, h.Login).Methods(http.MethodPost)

	return router
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var p LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	a, err := h.r.AdminRepository().FindByEmail(r.Context(), p.Email)
	if err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	if err := a.Authenticate(p.Password); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	set, err := h.r.JSONWebKeySetRepository().First(r.Context())
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = a.Email
	claims["kid"] = set.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(set.Private)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	resp := &LoginResponse{
		Token: tokenString,
	}
	h.r.Writer().Write(w, r, resp)
}

func (h *Handler) Authorize(claim jwt.Claims) (crypto.PublicKey, error) {
	kid, ok := claim.(jwt.MapClaims)["kid"]
	if !ok {
		h.r.Logger().Infoln("kid is not found in claim")
		return nil, errors.New("invalid claim")
	}
	set, err := h.r.JSONWebKeySetRepository().Find(context.Background(), fmt.Sprint(kid))
	if err != nil {
		h.r.Logger().Warnln(err)
		return nil, errors.New("invalid claim")
	}
	return set.Public, nil
}

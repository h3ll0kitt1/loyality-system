package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
	"github.com/h3ll0kitt1/loyality-system/internal/repository"
	"github.com/h3ll0kitt1/loyality-system/internal/service"
)

type Handlers struct {
	service Services
	log     *zap.SugaredLogger
}

func NewHandlers(s Services, log *zap.SugaredLogger) Handlers {
	return Handlers{
		service: s,
		log:     log,
	}
}

type Services interface {
	// auth
	CreateUser(ctx context.Context, credentials domain.Credentials) error
	AuthUser(ctx context.Context, credentials domain.Credentials) (string, error)

	// order
	CheckOrderIsNotDuplicated(ctx context.Context, username string, orderID uint32) (bool, error)
	LoadOrderInfo(ctx context.Context, username string, orderID uint32) error
	GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error)

	// // bonus
	GetBonusInfoForUser(ctx context.Context, username string) (domain.BonusInfo, error)
	WithdrawBonusForOrder(ctx context.Context, username string, orderID uint32, sum int64) error
	GetBonusOperationsForUser(ctx context.Context, username string) ([]domain.WithdrawInfo, error)
}

// Регистрация пользователя
func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var creds domain.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		h.writeResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.CreateUser(r.Context(), creds)
	switch {
	case errors.Is(err, repository.ErrUserAlreadyExists):
		h.writeResponse(w, http.StatusConflict, err)
		return
	case err != nil:
		h.writeResponse(w, http.StatusInternalServerError, err)
		return
	}

	authToken, err := h.service.AuthUser(r.Context(), creds)
	if err != nil {
		h.writeResponse(w, http.StatusInternalServerError, err)
		return
	}

	cookie := http.Cookie{
		Name:     "gohermart",
		Value:    authToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	h.writeResponse(w, http.StatusOK, "new user is registered")
}

func (h *Handlers) AuthUser(w http.ResponseWriter, r *http.Request) {

	var creds domain.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		h.writeResponse(w, http.StatusBadRequest, err)
		return
	}

	authToken, err := h.service.AuthUser(r.Context(), creds)
	switch {
	case errors.Is(err, repository.ErrUserNotExists):
		h.writeResponse(w, http.StatusUnauthorized, err)
		return
	case errors.Is(err, service.ErrWrongCredentials):
		h.writeResponse(w, http.StatusUnauthorized, err)
		return
	case err != nil:
		h.writeResponse(w, http.StatusInternalServerError, err)
		return
	}

	cookie := http.Cookie{
		Name:     "gohermart",
		Value:    authToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	h.writeResponse(w, http.StatusOK, "user is authorized")
}

func (h *Handlers) ErrorNotFound(w http.ResponseWriter, r *http.Request) {
	h.writeResponse(w, http.StatusNotFound, fmt.Errorf("page not found"))
}

func (h *Handlers) writeResponse(w http.ResponseWriter, status int, data interface{}) {

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if x, ok := data.(error); ok {
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: x.Error()}); err != nil {
			h.log.Infow("info",
				"write response failed:", err,
			)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Infow("info",
			"write response failed:", err,
		)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

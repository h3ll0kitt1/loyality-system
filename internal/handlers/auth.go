package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

type Handlers struct {
	service Services
}

func NewHandlers(s Services) Handlers {
	return Handlers{
		service: s,
	}
}

type Services interface {
	// auth
	CheckUserExists(ctx context.Context, username string) bool
	CreateUser(ctx context.Context, credentials domain.Credentials) error
	AuthUser(ctx context.Context, credentials domain.Credentials) (string, error)

	// bonus
	GetBonusInfoForUser(ctx context.Context, username string) (domain.BonusInfo, error)
	WithdrawBonusForOrder(ctx context.Context, username string, orderID uint64) bool
	GetBonusOperationsForUser(ctx context.Context, username string) ([]domain.WithdrawInfo, error)

	// order
	CheckOrderIsNotDuplicated(ctx context.Context, username string, orderID uint64) bool
	CheckOrderIsNotExistsForAnotherUser(ctx context.Context, username string, orderID uint64) bool
	ValidateWithLuhn(orderID uint64) bool
	LoadOrderInfo(ctx context.Context, username string, orderID uint64) error
	GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error)
}

// Регистрация пользователя
func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var input domain.Credentials
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := h.service.CheckUserExists(r.Context(), input.Username)
	if !ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = h.service.CreateUser(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authToken, err := h.service.AuthUser(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) AuthUser(w http.ResponseWriter, r *http.Request) {

	var input domain.Credentials
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authToken, err := h.service.AuthUser(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ErrorNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

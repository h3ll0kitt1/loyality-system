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
	"github.com/h3ll0kitt1/loyality-system/internal/utils"
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
	InsertOrderInfo(ctx context.Context, username string, orderID uint32) (bool, error)
	UpdateOrderInfo(ctx context.Context, order domain.OrderInfoRequest) error
	GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error)
	GetOrdersForUpdate(ctx context.Context, limit int32) ([]domain.OrderInfo, error)

	// bonus
	GetBonusInfoForUser(ctx context.Context, username string) (domain.BonusInfo, error)
	WithdrawBonusForOrder(ctx context.Context, username string, orderID uint32, sum int64) error
	GetBonusOperationsForUser(ctx context.Context, username string) ([]domain.WithdrawInfo, error)

	// updater

}

// Регистрация пользователя
func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var creds domain.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.service.CreateUser(r.Context(), creds)
	switch {
	case errors.Is(err, repository.ErrUserAlreadyExists):
		utils.WriteResponse(w, http.StatusConflict, err)
		return
	case err != nil:
		utils.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	authToken, err := h.service.AuthUser(r.Context(), creds)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, err)
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
	utils.WriteResponse(w, http.StatusOK, "new user is registered")
}

func (h *Handlers) AuthUser(w http.ResponseWriter, r *http.Request) {

	var creds domain.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	authToken, err := h.service.AuthUser(r.Context(), creds)
	switch {
	case errors.Is(err, repository.ErrUserNotExists):
		utils.WriteResponse(w, http.StatusUnauthorized, err)
		return
	case errors.Is(err, service.ErrWrongCredentials):
		utils.WriteResponse(w, http.StatusUnauthorized, err)
		return
	case err != nil:
		utils.WriteResponse(w, http.StatusInternalServerError, err)
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
	utils.WriteResponse(w, http.StatusOK, "user is authorized")
}

func (h *Handlers) ErrorNotFound(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, http.StatusNotFound, fmt.Errorf("page not found"))
}

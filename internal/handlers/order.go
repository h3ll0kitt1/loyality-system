package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/h3ll0kitt1/loyality-system/internal/crypto/validator"
	"github.com/h3ll0kitt1/loyality-system/internal/domain"
	"github.com/h3ll0kitt1/loyality-system/internal/repository"
	"github.com/h3ll0kitt1/loyality-system/internal/utils"
)

// загрузка пользователем номера заказа для расчёта
func (h *Handlers) LoadOrder(w http.ResponseWriter, r *http.Request) {

	login := r.Context().Value(domain.CtxLoginKey{})
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	orderID := string(body)
	_, err = strconv.ParseUint(orderID, 10, 64)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	ok := validator.LuhnAlgorithm(orderID)
	if !ok {
		utils.WriteResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	ok, err = h.service.InsertOrderInfo(r.Context(), login.(string), orderID)
	switch {
	case errors.Is(err, repository.ErrOrderAlreadyExistsForOtherUser):
		utils.WriteResponse(w, http.StatusConflict, err)
		return
	case err != nil:
		utils.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	if !ok {
		utils.WriteResponse(w, http.StatusOK, "order has been already registered for this user")
		return
	}
	utils.WriteResponse(w, http.StatusAccepted, "new order is registered")
}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
func (h *Handlers) GetOrdersInfo(w http.ResponseWriter, r *http.Request) {

	login := r.Context().Value(domain.CtxLoginKey{})

	ordersInfoHistory, err := h.service.GetOrdersInfoForUser(r.Context(), login.(string))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	if len(ordersInfoHistory) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(ordersInfoHistory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

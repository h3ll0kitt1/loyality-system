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
)

// загрузка пользователем номера заказа для расчёта
func (h *Handlers) LoadOrder(w http.ResponseWriter, r *http.Request) {

	login := r.Context().Value(domain.CtxLoginKey{})
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeResponse(w, http.StatusInternalServerError, err)
		return
	}

	orderID, err := strconv.ParseUint(string(body), 10, 64)
	if err != nil {
		h.writeResponse(w, http.StatusBadRequest, err)
		return
	}

	ok := validator.LuhnAlgorithm(uint32(orderID))
	if !ok {
		h.writeResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	ok, err = h.service.LoadOrderInfo(r.Context(), login.(string), uint32(orderID))
	switch {
	case errors.Is(err, repository.ErrOrderAlreadyExistsForOtherUser):
		h.writeResponse(w, http.StatusConflict, err)
		return
	case err != nil:
		h.writeResponse(w, http.StatusInternalServerError, err)
		return
	}

	if !ok {
		h.writeResponse(w, http.StatusOK, "order has been already registered for this user")
		return
	}
	h.writeResponse(w, http.StatusAccepted, "new order is registered")
}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
func (h *Handlers) GetOrdersInfo(w http.ResponseWriter, r *http.Request) {

	login := r.Context().Value(domain.CtxLoginKey{})

	ordersInfoHistory, err := h.service.GetOrdersInfoForUser(r.Context(), login.(string))
	if err != nil {
		h.writeResponse(w, http.StatusInternalServerError, err)
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

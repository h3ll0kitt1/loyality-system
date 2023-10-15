package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// загрузка пользователем номера заказа для расчёта
func (h *Handlers) LoadOrder(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orderID, err := strconv.ParseUint(string(body), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := h.service.ValidateWithLuhn(orderID)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ok = h.service.CheckOrderIsNotExistsForAnotherUser(r.Context(), username.(string), orderID)
	if !ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	ok = h.service.CheckOrderIsNotDuplicated(r.Context(), username.(string), orderID)
	if !ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	err = h.service.LoadOrderInfo(r.Context(), username.(string), orderID)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
func (h *Handlers) GetOrdersInfo(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username")

	ordersInfoHistory, err := h.service.GetOrdersInfoForUser(r.Context(), username.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

// получение текущего баланса счёта баллов лояльности пользователя
func (h *Handlers) GetBonusInfo(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username")
	bonusInfo, err := h.service.GetBonusInfoForUser(r.Context(), username.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(bonusInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *Handlers) WithdrawBonus(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username")

	var input domain.WithdrawInfo
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := h.service.ValidateWithLuhn(input.OrderID)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ok = h.service.WithdrawBonusForOrder(r.Context(), username.(string), input.OrderID)
	if !ok {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *Handlers) GetBonusOperationsInfo(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username")
	bonusInfoHistory, err := h.service.GetBonusOperationsForUser(r.Context(), username.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(bonusInfoHistory) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(bonusInfoHistory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

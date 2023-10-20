package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/h3ll0kitt1/loyality-system/internal/crypto/validator"
	"github.com/h3ll0kitt1/loyality-system/internal/domain"
	"github.com/h3ll0kitt1/loyality-system/internal/repository"
	"github.com/h3ll0kitt1/loyality-system/internal/utils"
)

// получение текущего баланса счёта баллов лояльности пользователя
func (h *Handlers) GetBonusInfo(w http.ResponseWriter, r *http.Request) {

	login := r.Context().Value(domain.CtxLoginKey{})
	bonusInfo, err := h.service.GetBonusInfoForUser(r.Context(), login.(string))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, err)
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

	login := r.Context().Value(domain.CtxLoginKey{})

	var withdraw domain.WithdrawInfo
	err := json.NewDecoder(r.Body).Decode(&withdraw)
	if err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, err)
		return
	}

	ok := validator.LuhnAlgorithm(withdraw.OrderID)
	if !ok {
		utils.WriteResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = h.service.WithdrawBonusForOrder(r.Context(), login.(string), withdraw.OrderID, withdraw.Sum)
	switch {
	case errors.Is(err, repository.ErrNotEnoughBonus):
		utils.WriteResponse(w, http.StatusPaymentRequired, err)
		return
	case err != nil:
		utils.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(w, http.StatusOK, "successful payment with bonus is done")
}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *Handlers) GetBonusOperationsInfo(w http.ResponseWriter, r *http.Request) {

	login := r.Context().Value(domain.CtxLoginKey{})
	bonusInfoHistory, err := h.service.GetBonusOperationsForUser(r.Context(), login.(string))
	if err != nil {
		utils.WriteResponse(w, http.StatusInternalServerError, err)
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

package domain

import (
	"time"
)

type BonusInfo struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawInfo struct {
	OrderID     uint64    `json:"order"`
	Sum         uint64    `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

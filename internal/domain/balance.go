package domain

import (
	"time"
)

type BonusInfo struct {
	Current  int64 `json:"current"`
	Withdraw int64 `json:"withdraw"`
}

type WithdrawInfo struct {
	OrderID     string    `json:"order" db:"order_id"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

package domain

import (
	"time"
)

type BonusInfo struct {
	Current  int64 `json:"current"`
	Withdraw int64 `json:"withdraw"`
}

type WithdrawInfo struct {
	OrderID     uint32    `json:"order" db:"id"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

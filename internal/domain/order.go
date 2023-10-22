package domain

import (
	"time"
)

type OrderInfo struct {
	Number     string    `json:"number" db:"order_id"`
	Status     string    `json:"status"`
	Accrual    int64     `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}

type OrderInfoRequest struct {
	Order   string `json:"order_id" db:"order_id"`
	Status  string `json:"status"`
	Accrual int64  `json:"accrual,omitempty"`
}

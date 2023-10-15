package domain

import (
	"time"
)

type OrderInfo struct {
	Number     uint64    `json:"number"`
	Status     string    `json:"status"`
	Accrual    uint64    `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

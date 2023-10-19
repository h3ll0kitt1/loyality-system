package domain

import (
	"time"
)

type OrderInfo struct {
	Number     uint32    `json:"number" db:"id"`
	Status     string    `json:"status"`
	Accrual    int64     `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}

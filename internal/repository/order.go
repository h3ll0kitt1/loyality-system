package repository

import (
	"context"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

func (r *RepositorySQL) CheckOrderIsNotDuplicated(ctx context.Context, username string, orderID uint64) bool {
	return true
}

func (r *RepositorySQL) CheckOrderIsNotExistsForAnotherUser(ctx context.Context, username string, orderID uint64) bool {
	return true
}

func (r *RepositorySQL) LoadOrderInfo(ctx context.Context, username string, orderID uint64) error {
	return nil
}

func (r *RepositorySQL) GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error) {
	return []domain.OrderInfo{}, nil
}

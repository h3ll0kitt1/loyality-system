package repository

import (
	"context"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

func (r *RepositorySQL) GetBonusInfoForUser(ctx context.Context, username string) (domain.BonusInfo, error) {
	return domain.BonusInfo{}, nil
}

func (r *RepositorySQL) WithdrawBonusForOrder(ctx context.Context, username string, orderID uint64) bool {
	return true
}

func (r *RepositorySQL) GetBonusOperationsForUser(ctx context.Context, username string) ([]domain.WithdrawInfo, error) {
	return []domain.WithdrawInfo{}, nil
}

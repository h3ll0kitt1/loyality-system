package service

import (
	"context"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

func (s *Service) GetBonusInfoForUser(ctx context.Context, username string) (domain.BonusInfo, error) {
	return s.repo.GetBonusInfoForUser(ctx, username)
}

func (s *Service) WithdrawBonusForOrder(ctx context.Context, username string, orderID uint32, sum int64) error {
	return s.repo.WithdrawBonusForOrder(ctx, username, orderID, sum)
}

func (s *Service) GetBonusOperationsForUser(ctx context.Context, username string) ([]domain.WithdrawInfo, error) {
	return s.repo.GetBonusOperationsForUser(ctx, username)
}

package service

import (
	"context"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

func (s *Service) CheckOrderIsNotDuplicated(ctx context.Context, username string, orderID uint32) (bool, error) {
	ok, err := s.repo.CheckOrderIsNotExistsForOtherUser(ctx, username, orderID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	ok, err = s.repo.CheckOrderIsNotExistsForThisUser(ctx, username, orderID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (s *Service) LoadOrderInfo(ctx context.Context, username string, orderID uint32) error {
	return s.repo.LoadOrderInfo(ctx, username, orderID)
}

func (s *Service) GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error) {
	return s.repo.GetOrdersInfoForUser(ctx, username)
}

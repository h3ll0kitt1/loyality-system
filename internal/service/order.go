package service

import (
	"context"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

func (s *Service) InsertOrderInfo(ctx context.Context, username string, orderID string) (bool, error) {
	return s.repo.InsertOrderInfo(ctx, username, orderID)
}

func (s *Service) UpdateOrderInfo(ctx context.Context, order domain.OrderInfoRequest) error {
	if order.Status == "INVALID" || order.Status == "PROCESSED" {
		return s.repo.UpdateOrderInfo(ctx, order)
	}
	return nil
}

func (s *Service) GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error) {
	return s.repo.GetOrdersInfoForUser(ctx, username)
}

func (s *Service) GetOrdersForUpdate(ctx context.Context, limit int32) ([]domain.OrderInfo, error) {
	return s.repo.GetOrdersForUpdate(ctx, limit)
}

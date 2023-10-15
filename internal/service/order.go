package service

import (
	"context"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

func (s *Service) CheckOrderIsNotDuplicated(ctx context.Context, username string, orderID uint64) bool {
	return s.repo.CheckOrderIsNotDuplicated(ctx, username, orderID)
}

func (s *Service) CheckOrderIsNotExistsForAnotherUser(ctx context.Context, username string, orderID uint64) bool {
	return s.repo.CheckOrderIsNotExistsForAnotherUser(ctx, username, orderID)
}

func (s *Service) ValidateWithLuhn(orderID uint64) bool {
	// implement luhn algorithm check
	return true
}

func (s *Service) LoadOrderInfo(ctx context.Context, username string, orderID uint64) error {

	// requestURL := fmt.Sprintf("http://%s/api/orders/%d", s.Addr, orderID)
	// req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	// if err != nil {
	// 	return err
	// }

	// make Request to GET /api/orders/{number} - чтобы отдать его системе обработки заказов
	// обработать запрос ? - если в ответе вернулись INVALID и PROCESSED значит все окей и помечаем в бд статус заказов
	// по идее если заказы помечены другими статусами то в надо в какой-тотгорутине периодически опрашивать текущие статусы и помечать бонусы
	return nil
}

func (s *Service) GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error) {
	return s.repo.GetOrdersInfoForUser(ctx, username)
}

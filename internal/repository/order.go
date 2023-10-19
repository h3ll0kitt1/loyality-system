package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

var (
	ErrOrderAlreadyExistsForOtherUser = fmt.Errorf("order has been already registered by another user")
)

func (r *RepositorySQL) CheckOrderIsNotExistsForOtherUser(ctx context.Context, username string, orderID uint32) (bool, error) {

	query := ` 	SELECT id FROM orders
   				WHERE id = $1 AND username <> $2`

	var result uint32

	err := r.db.GetContext(ctx, &result, query, orderID, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return true, nil
	case err != nil:
		return false, fmt.Errorf("repository: check order is not duplicated for another user failed: %w", err)
	}
	return false, ErrOrderAlreadyExistsForOtherUser
}

func (r *RepositorySQL) CheckOrderIsNotExistsForThisUser(ctx context.Context, username string, orderID uint32) (bool, error) {

	query := ` 	SELECT id FROM orders
   				WHERE id = $1 AND username = $2`

	var result uint32

	err := r.db.GetContext(ctx, &result, query, orderID, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return true, nil
	case err != nil:
		return false, fmt.Errorf("repository: check order is not duplicated for this user failed: %w", err)
	}
	return false, nil
}

func (r *RepositorySQL) LoadOrderInfo(ctx context.Context, username string, orderID uint32) error {

	query := `	INSERT INTO orders (id, username, accrual, status, uploaded_at)
				VALUES ($1, $2, -1, 'NEW', NOW())
				ON CONFLICT (id) DO NOTHING
				RETURNING TRUE`

	var ok bool

	err := r.db.GetContext(ctx, &ok, query, orderID, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("repository: load new order failed: %w", err)
	}
	return nil
}

func (r *RepositorySQL) GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error) {

	var orders []domain.OrderInfo

	query := `	SELECT id, status, accrual, uploaded_at FROM orders
   				WHERE username = $1`

	if err := r.db.SelectContext(ctx, &orders, query, username); err != nil {
		return nil, fmt.Errorf("repository: get orders for user failed: %w", err)
	}
	return orders, nil
}

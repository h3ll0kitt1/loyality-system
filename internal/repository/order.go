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

func (r *RepositorySQL) LoadOrderInfo(ctx context.Context, username string, orderID uint32) (bool, error) {

	tx, err := r.db.Beginx()
	if err != nil {
		return false, fmt.Errorf("repository: load order info failed: %w", err)
	}
	defer tx.Rollback()

	err = r.orderIsNotExistsForOtherUser(ctx, tx, username, orderID)
	if err != nil {
		return false, err
	}

	ok, err := r.orderIsNotExistsForThisUser(ctx, tx, username, orderID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	if err = r.insertOrderInfo(ctx, tx, username, orderID); err != nil {
		return false, err
	}

	if err = tx.Commit(); err != nil {
		return false, fmt.Errorf("repository: load order info failed: %w", err)
	}
	return true, nil
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

func (r *RepositorySQL) orderIsNotExistsForOtherUser(ctx context.Context, q q, username string, orderID uint32) error {

	query := ` 	SELECT id FROM orders
   				WHERE id = $1 AND username <> $2`

	var result uint32

	err := q.GetContext(ctx, &result, query, orderID, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("repository: check order is not duplicated for another user failed: %w", err)
	}
	return ErrOrderAlreadyExistsForOtherUser
}

func (r *RepositorySQL) orderIsNotExistsForThisUser(ctx context.Context, q q, username string, orderID uint32) (bool, error) {

	query := ` 	SELECT id FROM orders
   				WHERE id = $1 AND username = $2`

	var result uint32

	err := q.GetContext(ctx, &result, query, orderID, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return true, nil
	case err != nil:
		return false, fmt.Errorf("repository: check order is not duplicated for this user failed: %w", err)
	}
	return false, nil
}

func (r *RepositorySQL) insertOrderInfo(ctx context.Context, q q, username string, orderID uint32) error {

	query := `	INSERT INTO orders (id, username, accrual, status, uploaded_at)
				VALUES ($1, $2, -1, 'NEW', NOW())
				ON CONFLICT (id) DO NOTHING
				RETURNING TRUE`

	var ok bool

	err := q.GetContext(ctx, &ok, query, orderID, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("repository: load new order failed: %w", err)
	}
	return nil
}

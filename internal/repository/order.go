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

func (r *RepositorySQL) InsertOrderInfo(ctx context.Context, username string, orderID string) (bool, error) {

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

func (r *RepositorySQL) UpdateOrderInfo(ctx context.Context, order domain.OrderInfoRequest) error {

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("repository: load order info failed: %w", err)
	}
	defer tx.Rollback()

	query := `	UPDATE orders
				SET status = '$1', accrual = $2
				WHERE order_id = $3
				RETURNING TRUE`

	var ok bool

	err = tx.GetContext(ctx, &ok, query, order.Status, order.Accrual, order.Order)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("repository: update order info (update status, accrual) failed: %w", err)
	}

	if order.Status == "PROCESSED" {

		query := `	SELECT username FROM orders
   					WHERE order_id = $1`

		var username string
		err := tx.GetContext(ctx, &username, query, order.Order)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrUserNotExists
		case err != nil:
			return fmt.Errorf("repository: update order info (select username) failed: %w", err)
		}

		query = `	UPDATE bonus
					SET current = current + $1
					WHERE username = $2
					RETURNING TRUE`
		var ok bool
		err = tx.GetContext(ctx, &ok, query, order.Accrual, username)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		case err != nil:
			return fmt.Errorf("repository: update order info (update bonus) failed: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("repository: update order info failed: %w", err)
	}

	return nil
}

// query = `	CREATE TABLE IF NOT EXISTS orders (
// 			order_id VARCHAR(32) PRIMARY KEY,
// 			username VARCHAR(255) NOT NULL REFERENCES users(username),
// 			status VARCHAR(12) NOT NULL,
// 			accrual INTEGER DEFAULT -1,
// 			uploaded_at  TIMESTAMP NOT NULL)`
// _, err = tx.ExecContext(ctx, query)
// if err != nil {
// 	return nil, fmt.Errorf("create orders table failed: %w", err)
// }

// query = `	CREATE TABLE IF NOT EXISTS bonus (
// 			username VARCHAR(255) primary key REFERENCES users(username),
// 			current BIGINT DEFAULT 0,
// 			withdraw BIGINT DEFAULT 0)`

func (r *RepositorySQL) GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error) {

	var orders []domain.OrderInfo

	query := `	SELECT order_id, status, accrual, uploaded_at FROM orders
   				WHERE username = $1`

	if err := r.db.SelectContext(ctx, &orders, query, username); err != nil {
		return nil, fmt.Errorf("repository: get orders for user failed: %w", err)
	}
	return orders, nil
}

func (r *RepositorySQL) GetOrdersForUpdate(ctx context.Context, limit int32) ([]domain.OrderInfo, error) {

	var orders []domain.OrderInfo

	query := `	SELECT order_id, status FROM orders
   				WHERE status = 'NEW' OR status = 'PROCESSING'
   				LIMIT $1 `

	if err := r.db.SelectContext(ctx, &orders, query, limit); err != nil {
		return nil, fmt.Errorf("repository: get orders for update failed: %w", err)
	}
	return orders, nil
}

func (r *RepositorySQL) orderIsNotExistsForOtherUser(ctx context.Context, q q, username string, orderID string) error {

	query := ` 	SELECT order_id FROM orders
   				WHERE order_id = $1 AND username <> $2`

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

func (r *RepositorySQL) orderIsNotExistsForThisUser(ctx context.Context, q q, username string, orderID string) (bool, error) {

	query := ` 	SELECT order_id FROM orders
   				WHERE order_id = $1 AND username = $2`

	var result string

	err := q.GetContext(ctx, &result, query, orderID, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return true, nil
	case err != nil:
		return false, fmt.Errorf("repository: check order is not duplicated for this user failed: %w", err)
	}
	return false, nil
}

func (r *RepositorySQL) insertOrderInfo(ctx context.Context, q q, username string, orderID string) error {

	query := `	INSERT INTO orders (order_id, username, accrual, status, uploaded_at)
				VALUES ($1, $2, -1, 'NEW', NOW())
				ON CONFLICT (order_id) DO NOTHING
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

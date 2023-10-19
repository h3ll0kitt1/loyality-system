package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

var (
	ErrNotEnoughBonus = fmt.Errorf("not enough bonus for payment")
)

func (r *RepositorySQL) GetBonusInfoForUser(ctx context.Context, username string) (domain.BonusInfo, error) {

	query := ` 	SELECT current, withdraw FROM bonus
   				WHERE username = $1`

	var bonus domain.BonusInfo

	err := r.db.GetContext(ctx, &bonus, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.BonusInfo{}, ErrUserNotExists
	case err != nil:
		return domain.BonusInfo{}, fmt.Errorf("repository: get bonus info for user failed: %w", err)
	}
	return bonus, nil
}

func (r *RepositorySQL) WithdrawBonusForOrder(ctx context.Context, username string, orderID uint32, sum int64) error {

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("repository: withdraw bonus for order failed: %w", err)
	}
	defer tx.Rollback()

	if err := r.isEnoughBonus(ctx, tx, username, sum); err != nil {
		return err
	}

	if err = r.changeBonusBalance(ctx, tx, username, sum); err != nil {
		return fmt.Errorf("repository: change bonus balance failed: %w", err)
	}

	if err = r.updateBonusWithdrawInfo(ctx, tx, username, orderID, sum); err != nil {
		return fmt.Errorf("repository: update bonus withdrawals info failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("repository: withdraw bonus for order failed: %w", err)
	}
	return nil
}

func (r *RepositorySQL) GetBonusOperationsForUser(ctx context.Context, username string) ([]domain.WithdrawInfo, error) {

	query := `	SELECT id, sum, processed_at FROM withdraws
   				WHERE username = $1`

	var withdraws []domain.WithdrawInfo
	if err := r.db.SelectContext(ctx, &withdraws, query, username); err != nil {
		return nil, fmt.Errorf("repository: get bonus operations for user failed: %w", err)
	}
	return withdraws, nil
}

type q interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func (r *RepositorySQL) isEnoughBonus(ctx context.Context, q q, username string, sum int64) error {

	query := `	SELECT current FROM bonus
   				WHERE username = $1`

	var current int64
	err := q.GetContext(ctx, &current, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrUserAlreadyExists
	case err != nil:
		return fmt.Errorf("repository: check if enough bonus failed: %w", err)
	}

	if current < sum {
		return ErrNotEnoughBonus
	}
	return nil
}

func (r *RepositorySQL) changeBonusBalance(ctx context.Context, q q, username string, sum int64) error {

	query := `	UPDATE bonus
				SET withdraw = withdraw + $1, current = current - $1
				WHERE username = $2
				RETURNING TRUE`

	var ok bool

	err := q.GetContext(ctx, &ok, query, sum, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("repository: change bonus balance failed: %w", err)
	}
	return nil
}

func (r *RepositorySQL) updateBonusWithdrawInfo(ctx context.Context, q q, username string, orderID uint32, sum int64) error {

	query := `	INSERT INTO withdraws (id, username, sum, processed_at)
				VALUES ($1, $2, $3, NOW())
				ON CONFLICT (order_id) DO NOTHING
				RETURNING TRUE`

	var ok bool
	err := q.GetContext(ctx, &ok, query, orderID, username, sum)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("repository: update bonus withdrawn info failed: %w", err)
	}
	return nil
}

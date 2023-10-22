package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrUserNotExists     = fmt.Errorf("user doesn't exist")
)

type RepositorySQL struct {
	db  *sqlx.DB
	log *zap.SugaredLogger
}

func NewRepository(ctx context.Context, DatabaseDSN string, log *zap.SugaredLogger) (*RepositorySQL, error) {

	db, err := sqlx.ConnectContext(ctx, "pgx", DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("connection to database failed: %w", err)
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("initialization of database tables failed: %w", err)
	}
	defer tx.Rollback()

	query := `	CREATE TABLE IF NOT EXISTS users (
    			username VARCHAR(255) PRIMARY KEY,
    			password_hash VARCHAR(128) NOT NULL,
    			salt VARCHAR(32) NOT NULL)`
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("create users table failed: %w", err)
	}
	query = `	CREATE TABLE IF NOT EXISTS orders (
				order_id VARCHAR(32) PRIMARY KEY,
				username VARCHAR(255) NOT NULL REFERENCES users(username),
				status VARCHAR(12) NOT NULL,
				accrual INTEGER DEFAULT -1,
				uploaded_at  TIMESTAMP NOT NULL)`
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("create orders table failed: %w", err)
	}

	query = `	CREATE TABLE IF NOT EXISTS bonus (
				username VARCHAR(255) primary key REFERENCES users(username),
				current BIGINT DEFAULT 0,
				withdraw BIGINT DEFAULT 0)`
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("create bonus table failed: %w", err)
	}

	query = `	CREATE TABLE IF NOT EXISTS withdraws  (
				order_id VARCHAR(255) PRIMARY KEY,
				username VARCHAR(255) NOT NULL REFERENCES users(username),
				sum INTEGER NOT NULL,
				processed_at TIMESTAMP NOT NULL)`
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("create withdraws table failed: %w", err)
	}
	tx.Commit()

	return &RepositorySQL{
		db:  db,
		log: log,
	}, nil
}

func (r *RepositorySQL) CreateUser(ctx context.Context, username, passwordHash, salt string) error {

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("repository: create user failed: %w", err)
	}
	defer tx.Rollback()

	query := ` 	INSERT INTO users (username, password_hash, salt)
				VALUES ($1, $2, $3)
				ON CONFLICT (username) DO NOTHING 
				RETURNING TRUE`

	var ok bool
	err = tx.GetContext(ctx, &ok, query, username, passwordHash, salt)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrUserAlreadyExists
	case err != nil:
		return fmt.Errorf("repository: create user failed: %w", err)
	}

	query = ` 	INSERT INTO bonus (username)
				VALUES ($1)
				ON CONFLICT (username) DO NOTHING 
				RETURNING TRUE`

	err = tx.GetContext(ctx, &ok, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrUserAlreadyExists
	case err != nil:
		return fmt.Errorf("repository: create bonus failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("repository: create user failed: %w", err)
	}
	return nil
}

func (r *RepositorySQL) GetSaltForUser(ctx context.Context, username string) (string, error) {
	query := ` 	SELECT salt FROM users
   				WHERE username = $1`

	var salt string

	err := r.db.GetContext(ctx, &salt, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", ErrUserNotExists
	case err != nil:
		return "", fmt.Errorf("repository: get salt failed: %w", err)
	}
	return salt, nil
}

func (r *RepositorySQL) GetPasswordHashForUser(ctx context.Context, username string) (string, error) {
	query := ` 	SELECT password_hash FROM users
   				WHERE username = $1`

	var passwordHash string

	err := r.db.GetContext(ctx, &passwordHash, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", ErrUserNotExists
	case err != nil:
		return "", fmt.Errorf("repository: get password hash failed: %w", err)
	}
	return passwordHash, nil
}

package repository

import (
	"context"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

type RepositorySQL struct {
}

func NewRepository(DatabaseDSN string) (*RepositorySQL, error) {
	return &RepositorySQL{}, nil
}

func (r *RepositorySQL) CheckUserExists(ctx context.Context, username string) bool {
	return true
}

func (r *RepositorySQL) CreateUser(ctx context.Context, credentials domain.Credentials) error {
	return nil
}

func (r *RepositorySQL) GetSaltForUser(ctx context.Context, username string) string {
	return ""
}

func (r *RepositorySQL) GetPasswordHashForUser(ctx context.Context, username string) string {
	return ""
}

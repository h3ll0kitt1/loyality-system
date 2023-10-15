package service

import (
	"context"
	"errors"

	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

type Repository interface {

	// auth
	CheckUserExists(ctx context.Context, username string) bool
	CreateUser(ctx context.Context, credentials domain.Credentials) error
	GetSaltForUser(ctx context.Context, username string) string
	GetPasswordHashForUser(ctx context.Context, username string) string

	// balance
	GetBonusInfoForUser(ctx context.Context, username string) (domain.BonusInfo, error)
	WithdrawBonusForOrder(ctx context.Context, username string, orderID uint64) bool
	GetBonusOperationsForUser(ctx context.Context, username string) ([]domain.WithdrawInfo, error)

	// order
	CheckOrderIsNotDuplicated(ctx context.Context, username string, orderID uint64) bool
	CheckOrderIsNotExistsForAnotherUser(ctx context.Context, username string, orderID uint64) bool
	LoadOrderInfo(ctx context.Context, username string, orderID uint64) error
	GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error)
}

func (s *Service) CheckUserExists(ctx context.Context, username string) bool {
	return s.repo.CheckUserExists(ctx, username)
}

func (s *Service) CreateUser(ctx context.Context, credentials domain.Credentials) error {
	return s.repo.CreateUser(ctx, credentials)
}

func (s *Service) AuthUser(ctx context.Context, credentials domain.Credentials) (string, error) {

	if !s.repo.CheckUserExists(ctx, credentials.Username) {
		return "", errors.New("check username")
	}

	salt := s.repo.GetSaltForUser(ctx, credentials.Username)
	hashedPassword := s.repo.GetPasswordHashForUser(ctx, credentials.Username)

	if hashedPassword != s.hasher(credentials.Password, salt) {
		return "", errors.New("wrong credentials")
	}
	return s.createJWTAuthToken(credentials.Username), nil
}

func (s *Service) hasher(password string, salt string) string {
	// write impl of argon2 hasher
	return ""
}

func (s *Service) createJWTAuthToken(username string) string {
	// write impl JWT construction
	return ""
}

package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/h3ll0kitt1/loyality-system/internal/crypto/argon2"
	"github.com/h3ll0kitt1/loyality-system/internal/crypto/jwt"
	"github.com/h3ll0kitt1/loyality-system/internal/crypto/random"
	"github.com/h3ll0kitt1/loyality-system/internal/domain"
)

var (
	ErrWrongCredentials = fmt.Errorf("user has given wrong credentials")
)

type Service struct {
	repo Repository
	log  *zap.SugaredLogger
}

func NewService(repo Repository, log *zap.SugaredLogger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

type Repository interface {

	// auth
	CreateUser(ctx context.Context, username string, hashedPassword string, salt string) error
	GetSaltForUser(ctx context.Context, username string) (string, error)
	GetPasswordHashForUser(ctx context.Context, username string) (string, error)

	// order
	LoadOrderInfo(ctx context.Context, username string, orderID uint32) (bool, error)
	GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error)

	// // balance
	GetBonusInfoForUser(ctx context.Context, username string) (domain.BonusInfo, error)
	WithdrawBonusForOrder(ctx context.Context, username string, orderID uint32, sum int64) error
	GetBonusOperationsForUser(ctx context.Context, username string) ([]domain.WithdrawInfo, error)
}

func (s *Service) CreateUser(ctx context.Context, credentials domain.Credentials) error {

	salt := random.GenerateSalt()
	hashedPassword := argon2.GenerateHash(credentials.Password, salt)

	if err := s.repo.CreateUser(ctx, credentials.Login, hashedPassword, salt); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (s *Service) AuthUser(ctx context.Context, credentials domain.Credentials) (string, error) {

	salt, err := s.repo.GetSaltForUser(ctx, credentials.Login)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	hashedPassword, err := s.repo.GetPasswordHashForUser(ctx, credentials.Login)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if hashedPassword != argon2.GenerateHash(credentials.Password, salt) {
		return "", ErrWrongCredentials
	}

	return jwt.GenerateToken(credentials.Login), nil
}

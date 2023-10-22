package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/h3ll0kitt1/loyality-system/internal/config"
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
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config, log *zap.SugaredLogger) *Service {
	return &Service{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}

type Repository interface {

	// auth
	CreateUser(ctx context.Context, username string, hashedPassword string, salt string) error
	GetSaltForUser(ctx context.Context, username string) (string, error)
	GetPasswordHashForUser(ctx context.Context, username string) (string, error)

	// order
	InsertOrderInfo(ctx context.Context, username string, orderID string) (bool, error)
	UpdateOrderInfo(ctx context.Context, order domain.OrderInfoRequest) error

	GetOrdersInfoForUser(ctx context.Context, username string) ([]domain.OrderInfo, error)
	GetOrdersForUpdate(ctx context.Context, limit int32) ([]domain.OrderInfo, error)

	// // balance
	GetBonusInfoForUser(ctx context.Context, username string) (domain.BonusInfo, error)
	WithdrawBonusForOrder(ctx context.Context, username string, orderID string, sum int64) error
	GetBonusOperationsForUser(ctx context.Context, username string) ([]domain.WithdrawInfo, error)
}

func (s *Service) CreateUser(ctx context.Context, credentials domain.Credentials) error {

	salt, err := random.GenerateSalt()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

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

	token, err := jwt.GenerateToken(credentials.Login, s.cfg)
	if err != nil {
		return "", err
	}
	return token, nil
}

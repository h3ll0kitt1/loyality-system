package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/h3ll0kitt1/loyality-system/internal/crypto/random"
)

type Config struct {
	AccrualSystem struct {
		HostPort string
	}

	Gohermart struct {
		HostPort string
	}

	CheckInterval time.Duration
	DatabaseDSN   string

	JWT struct {
		TokenExpire time.Duration
		SecretKey   string
	}
}

func NewConfig() *Config {
	var cfg Config
	return &cfg
}

// надо возвращать ошибку если не заданы параметры которых нет по умолчанию ?
func (cfg *Config) Parse() error {

	var (
		flagCheckInterval int
		flagTokenExpire   int
		flagAccrualSystem string
		flagGohermart     string
		flagDatabaseDSN   string
	)

	flag.IntVar(&flagCheckInterval, "i", 30, "number of seconds to update order status")
	flag.IntVar(&flagTokenExpire, "e", 6, "number of minuts before JWT token expires for client")
	flag.StringVar(&flagAccrualSystem, "r", "localhost:8080", "address of system bonus calculations")
	flag.StringVar(&flagGohermart, "a", "localhost:9090", "address and port to run app")
	flag.StringVar(&flagDatabaseDSN, "d", "", "databaseDSN to connect to database")
	flag.Parse()

	envCheckInterval, err := strconv.Atoi(os.Getenv("CHECK_INTERVAL"))
	if err == nil {
		flagCheckInterval = envCheckInterval
	}

	envTokenExpire, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRE"))
	if err == nil {
		flagTokenExpire = envTokenExpire
	}

	if envAccrualSystem := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystem != "" {
		flagAccrualSystem = envAccrualSystem
	}

	if envGohermart := os.Getenv("RUN_ADDRESS"); envGohermart != "" {
		flagGohermart = envGohermart
	}

	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		flagDatabaseDSN = envDatabaseDSN
	}

	var envSecretKey string
	if envSecretKey = os.Getenv("SECRET_KEY"); envSecretKey == "" {
		secretKey, err := random.GenerateSecretKey()
		if err != nil {
			return fmt.Errorf("generating token without pre-set ENV failed, set token in ENV SECRET_KEY")
		}

		os.Setenv("SECRET_KEY", secretKey)
		envSecretKey = secretKey
	}

	cfg.AccrualSystem.HostPort = flagAccrualSystem
	cfg.Gohermart.HostPort = flagGohermart
	cfg.DatabaseDSN = flagDatabaseDSN
	cfg.CheckInterval = time.Duration(flagCheckInterval) * time.Second

	cfg.JWT.TokenExpire = time.Duration(flagTokenExpire) * time.Hour
	cfg.JWT.SecretKey = envSecretKey

	return nil
}

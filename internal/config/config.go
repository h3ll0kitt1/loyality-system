package config

import (
	"flag"
	"os"
)

type Config struct {
	SystemCalc string
	Server     struct {
		HostPort string
	}
	DatabaseDSN string
}

func NewConfig() *Config {
	var cfg Config
	return &cfg
}

// надо возвращать ошибку если не заданы параметры которых нет по умолчанию ?
func (cfg *Config) Parse() {

	var (
		flagAccrualSystem string
		flagHostPort      string
		flagDatabaseDSN   string
	)

	flag.StringVar(&flagAccrualSystem, "r", "", "address of system bonus calculations")
	flag.StringVar(&flagHostPort, "a", "localhost:8080", "address and port to run app")
	flag.StringVar(&flagDatabaseDSN, "d", "", "databaseDSN to connect to database")

	if envAccrualSystem := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystem != "" {
		flagAccrualSystem = envAccrualSystem
	}

	if envHostPort := os.Getenv("RUN_ADDRESS"); envHostPort != "" {
		flagHostPort = envHostPort
	}

	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		flagDatabaseDSN = envDatabaseDSN
	}

	cfg.SystemCalc = flagAccrualSystem
	cfg.Server.HostPort = flagHostPort
	cfg.DatabaseDSN = flagDatabaseDSN
}

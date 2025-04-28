package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBDSN     string
	Port      string
	LogLevel  string
	RateLimit int
}

func NewConfig() (*Config, error) {
	rateLimitStr := os.Getenv("RATE_LIMIT")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		return nil, err
	}
	return &Config{
		DBDSN:     os.Getenv("DB_DSN"),
		Port:      os.Getenv("PORT"),
		LogLevel:  os.Getenv("LOG_LEVEL"),
		RateLimit: rateLimit,
	}, nil

}

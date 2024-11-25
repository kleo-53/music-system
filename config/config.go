package config

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/kleo-53/music-system/pkg/logger"
)

type Config struct {
	Port     string
	DBURL    string
	LogLevel string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load("config.env"); err != nil {
		logger.Log().Warn(context.Background(), "No .env file found, using environment variables")
	}
	cfg := &Config{
		Port:     os.Getenv("HOST_PORT"),
		DBURL:    os.Getenv("DB_URL"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
	return cfg, nil
}

package config

import (
	"context"
	"os"

	"github.com/cyb3rkh4l1d/subsapi/internal/database"
	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Define configuration for the applications
// Определение конфигурации для приложений
type Config struct {
	Host     string
	LogLevel string
	GinMode  string
	DbConfig *database.Config
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/
// Function that loads configurations of the application from.env
// Функция, загружающая конфигурации приложения из файла .env
func LoadConfig(ctx context.Context, logger *logrus.Entry) *Config {
	err := godotenv.Load()
	cfg := &Config{

		Host:     getEnv("Host", ":8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		GinMode:  getEnv("GIN_MODE", "debug"),
		DbConfig: &database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgress"),
			Password: getEnv("DB_PASSWORD", "postgress"),
			DBName:   getEnv("DB_NAME", "subscriptions_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	if err != nil {
		logger.WithError(err).Warn(validations.ErrConfiLoadFailed)
	} else {
		logger.Info("Config Loaded Successfully")
	}

	return cfg
}

// function that gets enviroment variables
// Функция, которая получает переменные окружения
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

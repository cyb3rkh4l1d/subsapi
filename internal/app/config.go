package app

import (
	"log"
	"os"

	"github.com/cyb3rkh4l1d/subsapi/internal/database"
	"github.com/joho/godotenv"
)

// Define configuration for the applications
type Config struct {
	Port     string
	LogLevel string
	dbConfig *database.Config
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/
// Function that loads configurations of the application from.env
func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port: getEnv("APP_PORT", ":8080"),
		dbConfig: &database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "subscriptions_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	log.Printf("[+] Loaded config: %+v", cfg)
	return cfg
}

// function that gets enviroment variables
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

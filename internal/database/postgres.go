package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config defines the PostgreSQL database connection settings.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/

// NewPostgresConnection creates and returns a new GORM PostgreSQL connection.
// It builds the DSN from the provided configuration and verifies the connection.
func NewPostgresConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("[-] failed to connect to database: %w", err)
	}

	//Ping DB via SQL DB connection for verification
	_, err = db.DB()
	if err != nil {
		return nil, fmt.Errorf("[-] failed to get postgresSql db: %w", err)
	}

	return db, nil
}

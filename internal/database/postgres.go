package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config defines the PostgreSQL database connection settings.
// В конфигурации задаются параметры подключения к базе данных PostgreSQL.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Ensures only one instance of PgDriver exists throughout the application lifecycle.
// Гарантирует существование только одного экземпляра PgDriver на протяжении всего жизненного цикла приложения.
var (
	once             sync.Once
	PgDriverInstance *PgDriver
)

// PgDriver encapsulates database connectivity and migration capabilities.
// PgDriver инкапсулирует возможности подключения к базе данных и миграции.
type PgDriver struct {
	Gorm_DB     *gorm.DB
	Sql_DB      *sql.DB
	Db_Migrator gorm.Migrator
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/

// NewPostgresConnection creates and returns a new GORM PostgreSQL connection using the Singleton design pattern.
// It builds the DSN from the provided configuration and verifies the connection.
// Функция NewPostgresConnection создает и возвращает новое соединение GORM с PostgreSQL, используя шаблон проектирования Singleton.
// Он формирует DSN на основе предоставленной конфигурации и проверяет соединение.
func NewPostgresConnection(config *Config, dbLogger *logrus.Entry) *PgDriver {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			dbLogger.WithError(err).Fatal(validations.ErrDbConnectionFailed)

		}

		sqlDB, err := db.DB()
		if err != nil {
			dbLogger.WithError(err).Fatal(validations.ErrDbPingFailed)

		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		PgDriverInstance = &PgDriver{
			Gorm_DB:     db,
			Sql_DB:      sqlDB,
			Db_Migrator: db.Migrator(),
		}
	})
	return PgDriverInstance
}

// ClosePgDriverConnection safely closes the singleton PostgreSQL database connection pool.
// Функция ClosePgDriverConnection безопасно закрывает пул соединений с единственной базой данных PostgreSQL.
func ClosePgDriverConnection() {
	if PgDriverInstance.Sql_DB != nil {
		log.Println("closing database connection pool...")
		if err := PgDriverInstance.Sql_DB.Close(); err != nil {
			log.Println(validations.ErrDbCloseConnectionFailed)
		} else {
			log.Println("database connections closed successfully")
		}
	}
}

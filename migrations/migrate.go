package migrations

import (
	"github.com/cyb3rkh4l1d/subsapi/internal/database"
	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

// MigrateSubscriptions performs automatic database migration for the Subscription model.
// Uses goose to create or update the 'subscriptions' table schema based on the model.
// Returns an error if migration fails.
// MigrateSubscriptions выполняет автоматическую миграцию базы данных для модели Subscription.
// Использует goose для создания или обновления схемы таблицы 'subscriptions' на основе модели.
// Возвращает ошибку, если миграция не удалась.
func PostgreSQLMigrateSubscriptions(dbLogger *logrus.Entry) {
	if err := goose.SetDialect("postgres"); err != nil {
		dbLogger.WithError(err).Fatal(validations.ErrDbMigrationFailed)

	}
	if err := goose.Up(database.PgDriverInstance.Sql_DB, "migrations"); err != nil {
		dbLogger.WithError(err).Fatal(validations.ErrDbMigrationFailed)
	}

	dbLogger.Info("database migration successful.")
}

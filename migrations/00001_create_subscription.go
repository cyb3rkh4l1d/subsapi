package migrations

import (
	"context"
	"database/sql"

	"github.com/cyb3rkh4l1d/subsapi/internal/database"
	"github.com/cyb3rkh4l1d/subsapi/internal/models"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateSubscription, downCreateSubscription)
}

func upCreateSubscription(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	database.PgDriverInstance.Db_Migrator.CreateTable(&models.Subscription{})
	return nil
}

func downCreateSubscription(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	database.PgDriverInstance.Db_Migrator.DropTable(&models.Subscription{})
	return nil
}

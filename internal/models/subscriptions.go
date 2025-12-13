package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Subscription represents a subscription record in the database.
// Maps directly to the 'subscriptions' table in PostgreSQL with GORM annotations.
type Subscription struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ServiceName string     `gorm:"type:varchar(100);not null" json:"service_name" example:"Yandex Plus"`
	Price       int        `gorm:"not null" json:"price" example:"400"`
	UserID      string     `gorm:"type:uuid;not null" json:"user_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	StartDate   time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate     *time.Time `gorm:"type:date" json:"end_date,omitempty"`
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/
// MigrateSubscriptions performs automatic database migration for the Subscription model.
// Uses GORM's AutoMigrate to create or update the 'subscriptions' table schema based on the model.
// Returns an error if migration fails.
func MigrateSubscriptions(db *gorm.DB) error {
	if err := db.AutoMigrate(&Subscription{}); err != nil {
		return fmt.Errorf("failed to migrate subscriptions table: %w", err)
	}
	return nil
}

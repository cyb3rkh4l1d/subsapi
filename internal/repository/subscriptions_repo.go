package repository

import (
	"context"
	"errors"
	"time"

	"github.com/cyb3rkh4l1d/subsapi/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SubscriptionRepository manages CRUD operations for subscriptions.
// It uses GORM for database access and Logrus for logging.
type SubscriptionRepository struct {
	DB     *gorm.DB
	Logger *logrus.Entry
}

/*
.....................................................................

	Functions/Methods Definations

........................................................................
*/
// NewSubscriptionRepository initializes a new repository instance.
func NewSubscriptionRepository(db *gorm.DB, logger *logrus.Entry) *SubscriptionRepository {
	return &SubscriptionRepository{
		DB:     db,
		Logger: logger,
	}
}

// CreateSubscription inserts a new subscription into the database.
func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	r.Logger.WithFields(logrus.Fields{
		"user_id": sub.UserID,
		"service": sub.ServiceName,
	}).Info("Creating subscription")

	result := r.DB.WithContext(ctx).Create(sub)

	if result.Error != nil {
		r.Logger.WithError(result.Error).Error("[-] failed to create subscription")
	}

	return result.Error
}

// GetByID retrieves a subscription by its ID.
func (r *SubscriptionRepository) GetByID(ctx context.Context, id uint) (*models.Subscription, error) {
	var sub models.Subscription
	if err := r.DB.WithContext(ctx).First(&sub, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.Logger.Errorf("[-] GetByID error: %v", err)
		return nil, err
	}
	return &sub, nil
}

// List fetches all subscriptions.
func (r *SubscriptionRepository) List(ctx context.Context) ([]models.Subscription, error) {
	var subs []models.Subscription
	if err := r.DB.WithContext(ctx).Find(&subs).Error; err != nil {
		r.Logger.Errorf("[-] List error: %v", err)
		return nil, err
	}
	return subs, nil
}

// Update saves changes to a subscription.
func (r *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	if err := r.DB.WithContext(ctx).Save(sub).Error; err != nil {
		r.Logger.Errorf("[-] Update error: %v", err)
		return err
	}
	return nil
}

// Delete removes a subscription by ID.
func (r *SubscriptionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&models.Subscription{}, id).Error; err != nil {
		r.Logger.Errorf("[-] Delete error: %v", err)
		return err
	}
	return nil
}

// CalculateTotalCost calculates total subscription cost and count for a user.
// Optionally filters by service name and start date range.
func (r *SubscriptionRepository) CalculateTotalCost(
	ctx context.Context,
	periodStart, periodEnd time.Time,
	userID, serviceName string,
) (int64, int64, error) {
	var totalCost, count int64

	// Build base query and arguments
	query := `
    SELECT 
        COALESCE(SUM(price), 0),
        COUNT(*)
    FROM subscriptions
    WHERE user_id = ?
    `
	args := []interface{}{userID}

	// Optional service name filter
	if serviceName != "" {
		query += " AND service_name = ?"
		args = append(args, serviceName)
	}

	// Optional date filters on start_date
	if !periodStart.IsZero() {
		query += " AND start_date >= ?"
		args = append(args, periodStart)
	}
	if !periodEnd.IsZero() {
		query += " AND start_date <= ?"
		args = append(args, periodEnd)
	}

	// Use GORM's Row() and Scan
	row := r.DB.WithContext(ctx).Raw(query, args...).Row()
	err := row.Scan(&totalCost, &count)
	if err != nil {
		return 0, 0, err
	}

	return totalCost, count, nil
}

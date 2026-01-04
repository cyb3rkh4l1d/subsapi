package repository

import (
	"context"
	"errors"

	"github.com/cyb3rkh4l1d/subsapi/internal/models"
	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Repository defines data access operations for subscription management
// Репозиторий определяет операции доступа к данным для управления подписками
type Repository interface {
	CreateSubscription(ctx context.Context, sub *models.Subscription) error
	GetSubscriptionByID(ctx context.Context, id uint) (*models.Subscription, error)
	ListSubscription(ctx context.Context, req *models.ListSubscriptionRequest) (int64, []models.Subscription, error)
	UpdateSubscriptionByID(ctx context.Context, sub *models.Subscription) error
	DeleteSubscriptionByID(ctx context.Context, id uint) error
	FindSubscriptionsByUserIDandServiceName(ctx context.Context, userID string, serviceName string) ([]models.Subscription, error)
}

// SubscriptionRepository manages CRUD operations for subscriptions.
// It uses GORM for database access and Logrus for logging.
// SubscriptionRepository управляет операциями CRUD для подписок.
// Он использует GORM для доступа к базе данных и Logrus для ведения журналов.
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
// NewSubscriptionRepository инициализирует новый экземпляр репозитория.
func NewSubscriptionRepository(db *gorm.DB, logger *logrus.Entry) *SubscriptionRepository {
	return &SubscriptionRepository{
		DB:     db,
		Logger: logger,
	}
}

// CreateSubscription inserts a new subscription into the database.
// Функция CreateSubscription вставляет новую подписку в базу данных.
func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	result := r.DB.WithContext(ctx).Create(sub)

	if result.Error != nil {
		r.Logger.WithError(result.Error).Error(validations.ErrCreateSubscriptionFailed)
		return validations.ErrCreateSubscriptionFailed
	}

	r.Logger.Info("subscription has been created:", *sub)
	return nil
}

// GetSubscriptionByID retrieves a subscription by its ID.
// Функция GetBGetSubscriptionByIDyID извлекает подписку по ее идентификатору.
func (r *SubscriptionRepository) GetSubscriptionByID(ctx context.Context, id uint) (*models.Subscription, error) {
	var sub models.Subscription
	if err := r.DB.WithContext(ctx).First(&sub, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.Logger.WithError(err).Error(validations.ErrGetSubscriptionByIDFailed)
		return nil, validations.ErrGetSubscriptionByIDFailed
	}
	r.Logger.Infof("subscription %+v has been fetched successfully: ", sub.ID)
	return &sub, nil
}

// ListSubscription fetches all subscriptions.
// ListSubscription получает все подписки.
func (r *SubscriptionRepository) ListSubscription(ctx context.Context, req *models.ListSubscriptionRequest) (int64, []models.Subscription, error) {
	var total int64
	var subs []models.Subscription
	orderClause := req.SortBy + " " + req.Order

	// count all subscriptions
	// подсчитать все подписки
	if err := r.DB.WithContext(ctx).Model(&models.Subscription{}).Count(&total).Error; err != nil {
		r.Logger.WithError(err).Error(validations.ErrListSubscriptionFailed)
		return total, nil, validations.ErrListSubscriptionFailed
	}

	//retrieves user's subscriptions with filtering, pagination, and sorting
	//Получает подписки пользователей с фильтрацией, пагинацией и сортировкой.
	if err := r.DB.WithContext(ctx).Limit(req.Limit).Offset(req.Offset).Order(orderClause).Find(&subs).Error; err != nil {
		r.Logger.WithError(err).Error(validations.ErrListSubscriptionFailed)
		return total, nil, validations.ErrListSubscriptionFailed
	}
	return total, subs, nil
}

// UpdateSubscription updates given subscription by its ID
// Функция UpdateSubscription обновляет указанную подписку по ее идентификатору.
func (r *SubscriptionRepository) UpdateSubscriptionByID(ctx context.Context, sub *models.Subscription) error {
	if err := r.DB.WithContext(ctx).Save(sub).Error; err != nil {
		r.Logger.WithError(err).Error(validations.ErrUpdateSubscriptionFailed)
		return validations.ErrUpdateSubscriptionFailed
	}
	r.Logger.Infof("subscription %+v has been updated successfully: ", sub.ID)
	return nil
}

// DeleteSubscription removes a subscription by ID.
// Функция DeleteSubscription удаляет подписку по ID.
func (r *SubscriptionRepository) DeleteSubscriptionByID(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&models.Subscription{}, id).Error; err != nil {
		r.Logger.WithError(err).Error(validations.ErrDeleteSubscriptionFailed)
		return validations.ErrDeleteSubscriptionFailed
	}
	r.Logger.Infof("subscription %+v has been deleted: ", id)
	return nil
}

// FindSubscriptionsByUserIDandServiceName Get subscriptions filtered by user and service_name
// FindSubscriptionsByUserIDandServiceName Получает подписки, отфильтрованные по пользователю и имени сервиса.
func (r *SubscriptionRepository) FindSubscriptionsByUserIDandServiceName(
	ctx context.Context,
	userID string,
	serviceName string,
) ([]models.Subscription, error) {
	query := r.DB.WithContext(ctx).Model(&models.Subscription{}).
		Where("user_id = ? AND service_name = ?", userID, serviceName)

	var subscriptions []models.Subscription
	if err := query.Find(&subscriptions).Error; err != nil {
		r.Logger.WithError(err).Error(validations.ErrFindSubscriptionByPeriodFailed)
		return nil, err
	}

	r.Logger.Infof("subscriptions for user %+v has been fetched: %+v", userID, subscriptions)
	return subscriptions, nil
}

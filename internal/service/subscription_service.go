package service

import (
	"context"
	"time"

	"github.com/cyb3rkh4l1d/subsapi/internal/models"
	"github.com/cyb3rkh4l1d/subsapi/internal/repository"
	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/sirupsen/logrus"
)

// SubscriptionService manages business logic for subscriptions
// SubscriptionService управляет бизнес-логикой для подписок
type SubscriptionService struct {
	repo   repository.Repository
	Logger *logrus.Entry
}

// NewSubscriptionService creates a new subscription service
// NewSubscriptionService создает новую службу подписки
func NewSubscriptionService(repo repository.Repository, logger *logrus.Entry) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		Logger: logger,
	}
}

// CreateSubscription handles business logic for creating a subscription
// Функция CreateSubscription обрабатывает бизнес-логику создания подписки
func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {

	//validate userId
	//проверить UserID
	err := validations.ValidateUserID(req.UserID)
	if err != nil {
		return nil, err
	}
	// Parse start_date (MM-YYYY)
	//проверить start_date
	startDate, err := validations.ValidateStartDate(req.StartDate)
	if err != nil {
		return nil, err
	}

	// Parse optional end_date (MM-YYYY)
	//проверить end_date
	endDate, err := validations.ValidateEndDate(startDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	//Validate service_name
	//проверить service_name
	if err := validations.ValidateServiceName(req.ServiceName); err != nil {
		return nil, err
	}

	// Create a subscription object based on the request data
	// Создание объекта подписки на основе данных запроса
	sub := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	// Save to database
	//Сохранить в базу данных
	if err := s.repo.CreateSubscription(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// GetSubscription retrieves a subscription by ID
// Метод GetSubscription извлекает подписку по ID
func (s *SubscriptionService) GetSubscription(ctx context.Context, id uint) (*models.Subscription, error) {

	// Retrieve the subscription by ID from the repository
	// Получение подписки по ID из репозитория
	sub, err := s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, validations.ErrSubscriptionNotFound
	}
	return sub, nil
}

// ListSubscriptions retrieves user's subscriptions with filtering, pagination, and sorting
// ListSubscriptions извлекает подписки пользователя с фильтрацией, пагинацией и сортировкой.
func (s *SubscriptionService) ListSubscriptions(ctx context.Context, req *models.ListSubscriptionRequest) (int64, []models.Subscription, error) {

	// retrieves user's subscriptions
	//Получить подписки пользователей
	total, subs, err := s.repo.ListSubscription(ctx, req)
	if err != nil {
		return total, nil, err
	}

	return total, subs, nil
}

// UpdateSubscription handles business logic for updating a subscription
// Функция UpdateSubscription обрабатывает бизнес-логику обновления подписки
func (s *SubscriptionService) UpdateSubscriptionByID(ctx context.Context, id uint, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	// check if subscription exists
	// Проверить, существует ли подписка
	sub, err := s.GetSubscription(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update service name if provided
	// Обновите имя службы, если оно указано
	if req.ServiceName != "" {
		if err := validations.ValidateServiceName(req.ServiceName); err != nil {
			return nil, err
		}
		sub.ServiceName = req.ServiceName
	}

	//update startdate if provided.
	//Обновите дату начала, если она указана.
	if req.StartDate != "" {
		startDate, err := validations.ValidateStartDate(req.StartDate)
		if err != nil {
			return nil, err
		}
		sub.StartDate = startDate
	}
	//update price if provided.
	//Обновить цену, если она указана.
	if req.Price > 0 {
		sub.Price = req.Price
	}
	// Update or clear end date and enforce end_date >= start_date
	// Обновить или очистить конечную дату и установить значение end_date >= start_date

	if req.EndDate == "" {
		sub.EndDate = nil
	} else {
		endDate, err := validations.ValidateEndDate(sub.StartDate, req.EndDate)
		if err != nil {
			return nil, err
		}
		sub.EndDate = endDate
	}

	// Save updates to the database
	// Сохранение обновлений в базу данных
	if err := s.repo.UpdateSubscriptionByID(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// The GetUserSubscriptionSummary function calculates and returns subscription statistics for a user.
// Функция GetUserSubscriptionSummary вычисляет и возвращает статистику подписки для пользователя.
func (s *SubscriptionService) GetUserSubscriptionSummary(
	ctx context.Context,
	req *models.UserSubscriptionSummaryRequest,
) (int, int64, int, error) {

	var periodStart time.Time
	var periodEnd *time.Time

	//validate userId
	//проверить UserID
	err := validations.ValidateUserID(req.UserID)
	if err != nil {
		return 0, 0, 0, err
	}

	//Validate service_name
	//проверить service_name
	if err := validations.ValidateServiceName(req.ServiceName); err != nil {
		return 0, 0, 0, err
	}

	//Validate query "from"
	// if query "from" is empty, periodstart default to time.TIme{}, otherwise it validate the query "from" value.
	//проверить query "from"
	// Если значение параметра "from" в запросе пустое, periodstart по умолчанию равен time.TIme{}, в противном случае выполняется проверка значения параметра "from" в запросе.
	if req.From == "" {
		periodStart = time.Time{}
	} else {
		periodStart, err = validations.ValidateStartDate(req.From)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	//Validate query "to"
	//if no query "to" is given in the query, periodEnd default  to current time, otherwise it validate the query "to" value.
	//проверить query "to"
	//Если в запросе не указан параметр "to", periodEnd по умолчанию принимает текущее время, в противном случае выполняется проверка значения параметра "to" в запросе.
	if req.To == "" {
		now := time.Now()
		periodEnd = &now
	} else {
		periodEnd, err = validations.ValidateEndDate(periodStart, req.To)
		if err != nil {
			return 0, 0, 0, err
		}
	}

	// Get all subscriptions for user
	// Получить все подписки пользователя
	subscriptions, err := s.repo.FindSubscriptionsByUserIDandServiceName(ctx, req.UserID, req.ServiceName)
	if err != nil {
		return 0, 0, 0, err
	}

	// Calculate total cost and unique months for user's subscription
	// Рассчитать общую стоимость и количество уникальных месяцев подписки пользователя
	unitPrice, totalCost, totalUniqueMonths := CalculateSubscriptionMetrics(
		subscriptions,
		periodStart,
		*periodEnd,
	)

	s.Logger.Infof("subscription metrics: UserID: %+v, ServiceName: %+v, TotalMonths: %+v, TotalCost: %+v", req.UserID, req.ServiceName, totalUniqueMonths, totalCost)

	return unitPrice, totalCost, totalUniqueMonths, nil
}

// DeleteSubscription deletes a subscription by its ID
// Функция DeleteSubscription удаляет подписку по её ID
func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id uint) error {
	// Check if subscription exists before deleting
	// Перед удалением проверьте, существует ли подписка.
	sub, err := s.GetSubscription(ctx, id)
	if err != nil {
		return err
	}
	// Delete the subscription from the database
	// Удалить подписку из базы данных
	if err := s.repo.DeleteSubscriptionByID(ctx, sub.ID); err != nil {
		return err
	}

	return nil
}

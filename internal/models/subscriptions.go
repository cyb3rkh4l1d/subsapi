package models

import (
	"time"
)

// Subscription represents a subscription record in the database.
// Maps directly to the 'subscriptions' table in PostgreSQL with GORM annotations.
// Indexes: Primary key (ID), composite index on (UserID, ServiceName).
// Subscription представляет собой запись о подписке в базе данных.
// Сопоставляется напрямую с таблицей 'subscriptions' в PostgreSQL с использованием аннотаций GORM.
// Индексы: первичный ключ (ID), составной индекс по (UserID, ServiceName).
type Subscription struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      string     `gorm:"type:uuid;not null;index:idx_summary_service,priority:1" json:"user_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	ServiceName string     `gorm:"type:varchar(100);not null;index:idx_summary_service,priority:2" json:"service_name" example:"Yandex Plus"`
	Price       int        `gorm:"not null" json:"price" example:"400"`
	StartDate   time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate     *time.Time `gorm:"type:date" json:"end_date" binding:"omitempty"`
}

// @Description Defines the request body for creating a new subscription.
// Определяет тело запроса для создания новой подписки.
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required,max=15"`
	Price       int    `json:"price" binding:"required,gt=0"`
	UserID      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date,omitempty"`
}

// @Description Defines the request body for updating a subscription.
// Определяет тело запроса для обновления подписки.
type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"omitempty,max=15"`
	Price       int    `json:"price" binding:"omitempty,gt=0"`
	StartDate   string `json:"start_date" binding:"omitempty"`
	EndDate     string `json:"end_date" binding:"omitempty"`
}

// @Description Defines the API response structure for a subscription.
// Определяет структуру ответа API для подписки.
type SubscriptionResponse struct {
	ID          uint   `json:"service_id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

// @Description Defines the request query for fetching subscription summary of a user.
// Определяет запрос для получения сводной информации о подписке пользователя.
type UserSubscriptionSummaryRequest struct {
	UserID      string `form:"user_id" binding:"required,uuid"`
	ServiceName string `form:"service_name,omitempty" binding:"required"`
	From        string `form:"from,omitempty"`
	To          string `form:"to,omitempty"`
}

// @Description Defines the generic error
// Определяет общую ошибку
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// @Description Defines the structure of the API response for the /summary endpoint.
// Определяет структуру ответа API для конечной точки /summary.
type UserSubscriptionSummaryResponse struct {
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	UnitPrice   int    `json:"unit_price"`
	TotalMonths int    `json:"total_months"`
	TotalAmount int64  `json:"total_amount"`
}

// @Description Defines the request query for fetching subscriptions with pagination, sorting and ordering
// Определяет запрос для получения подписок с пагинацией, сортировкой и упорядочиванием.
type ListSubscriptionRequest struct {
	Limit  int    `form:"limit,default=10" json:"limit" binding:"omitempty,min=1,max=100"`                      // Max items to return
	Offset int    `form:"offset,default=0" json:"offset" binding:"omitempty,min=0"`                             // Items to skip
	SortBy string `form:"sort_by,default=id" binding:"oneof=id user_id service_name price start_date end_date"` // created_at, price, start_date
	Order  string `form:"order,default=desc" binding:"oneof=desc asc"`                                          // asc, desc
}

// @Description Defines the request query path processing subscription by ID
// Определяет подписку на обработку пути запроса по идентификатору.
type SubscriptionUriIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

// @Description Defines pagination metadata for response for ListSubscriptionResponse
// Определяет метаданные для пагинации в ответе для ListSubscriptionResponse
type PaginationMeta struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
	Total  int64  `json:"total"`
}

// @Description Defines the API response structure for a ListSubscriptionRequest.
// Определяет структуру ответа API для запроса ListSubscriptionRequest.
type ListSubscriptionsResponse struct {
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
	Meta          *PaginationMeta        `json:"meta"`
}

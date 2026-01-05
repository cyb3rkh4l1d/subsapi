package handlers

import (
	"context"
	"net/http"

	"github.com/cyb3rkh4l1d/subsapi/internal/models"
	"github.com/cyb3rkh4l1d/subsapi/internal/service"
	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SubscriptionHandler handles HTTP requests related to subscriptions.
// It contains shared context, logger, and repository dependencies.
// SubscriptionHandler обрабатывает HTTP-запросы, связанные с подписками.
// Он содержит зависимости от общего контекста, логгера и репозитория.
type SubscriptionHandler struct {
	ctx     context.Context
	Logger  *logrus.Entry
	service *service.SubscriptionService
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/

// NewSubscriptionHandlers creates and returns a SubscriptionHandler instance with
// With shared context, logger, and repository dependencies.
// NewSubscriptionHandlers создает и возвращает экземпляр SubscriptionHandler с
// С зависимостями от общего контекста, логгера и репозитория.
func NewSubscriptionHandlers(ctx context.Context, handlerLogger *logrus.Entry, serice *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{ctx: ctx, Logger: handlerLogger, service: serice}
}

// @tag.name Subscriptions
// @tag.description Subscription management endpoints

// CreateSubscription handles HTTP POST requests to create a new subscription.
// It validates input, parses dates, persists data, and returns the created record.
// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a subscription for a user
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.CreateSubscriptionRequest true "Subscription payload"
// @Success 201 {object} models.SubscriptionResponse
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {

	var req *models.CreateSubscriptionRequest

	// Bind and validate request payload
	//Привяжите и проверьте полезную нагрузку запроса.
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.WithError(err).Info(validations.ErrInvalidRequestInput)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: validations.ErrInvalidRequestInput.Error(), Details: err.Error()})
		return
	}

	h.Logger.Infof("creating subscription: ServiceName: %+v, UserID: %+v, Price: %+v,StartDate: %+v, EndDate: %+v", req.ServiceName, req.UserID, req.Price, req.StartDate, req.EndDate)

	//Process business logic for create subscription request
	//Обработка бизнес-логики для создания запроса на подписку
	sub, err := h.service.CreateSubscription(c.Request.Context(), req)

	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, FormatToSubscriptionResponse(sub))
}

// ListSubscriptions retrieves paginated subscriptions with optional sorting and filtering
// It converts internal date fields to MM-YYYY format and returns a paginated API response
// ListSubscriptions godoc
// @Summary List subscriptions with pagination
// @Description Retrieve paginated list of subscriptions with optional sorting
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items to return" default(10) minimum(1) maximum(100)
// @Param offset query int false "Number of items to skip" default(0) minimum(0)
// @Param sort_by query string false "Field to sort by" default(id) Enums(id, user_id, service_name, price, start_date, end_date)
// @Param order query string false "Sort order" default(desc) Enums(asc, desc)
// @Success 200 {object} models.ListSubscriptionsResponse
// @Failure 400 {object} models.ErrorResponse "Bad Request - Invalid query parameters"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {

	var req *models.ListSubscriptionRequest

	// Bind and validate request payload
	//Привяжите и проверьте полезную нагрузку запроса.
	if err := c.ShouldBindQuery(&req); err != nil {
		h.Logger.WithError(err).Info(validations.ErrInvalidRequestInput)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: validations.ErrInvalidRequestInput.Error(), Details: err.Error()})
		return
	}
	h.Logger.Infof("getting subscriptions:- Limit: %+v, Offset: %+v, SortBy: %+v, Order: %+v", req.Limit, req.Offset, req.SortBy, req.Order)

	//process business logic for ListSubscriptionRequest
	//Обработка бизнес-логики для ListSubscriptionRequest
	total, subs, err := h.service.ListSubscriptions(c.Request.Context(), req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Convert each subscription model to API response format
	// Преобразовать каждую модель подписки в формат ответа API
	formatedSubs := make([]models.SubscriptionResponse, len(subs))
	for i, sub := range subs {
		formatedSubs[i] = FormatToSubscriptionResponse(&sub)
	}

	// Create pagination metadata for the response
	// Создание метаданных для пагинации ответа
	paginationMeta := &models.PaginationMeta{Limit: req.Limit, Offset: req.Offset, SortBy: req.SortBy, Order: req.Order, Total: total}

	// Create a final response with subscription and pagination data
	// Создать окончательный ответ с данными о подписке и постраничной навигации
	res := &models.ListSubscriptionsResponse{Subscriptions: formatedSubs, Meta: paginationMeta}

	c.JSON(http.StatusOK, res)

}

// GetSubscription retrieves a single subscription by its ID.
// It validates the identifier and returns a formatted subscription response if found.
// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Retrieve a subscription using its ID
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID" minimum(1)
// @Success 200 {object} models.SubscriptionResponse
// @Failure 400 {object} models.ErrorResponse "Bad Request - Invalid subscription ID"
// @Failure 404 {object} models.ErrorResponse "Not Found - Subscription does not exist"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {

	var req *models.SubscriptionUriIDRequest

	// Bind and validate request payload
	//Привяжите и проверьте полезную нагрузку запроса.
	if err := c.ShouldBindUri(&req); err != nil {
		h.Logger.WithError(err).Info(validations.ErrInvalidRequestInput)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: validations.ErrInvalidRequestInput.Error(), Details: err.Error(),
		})
		return
	}

	h.Logger.Info("getting subscription by ID: ", req.ID)

	//process business logic for GetSubscriptionRequest
	//Обработка бизнес-логики для GetSubscription Request
	sub, err := h.service.GetSubscription(c.Request.Context(), req.ID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, FormatToSubscriptionResponse(sub))

}

// UpdateSubscription updates an existing subscription by ID.
// Only fields provided in the request are modified (partial update/PATCH-like),
// with validation applied to price and date formats.
// UpdateSubscription godoc
// @Summary Update subscription
// @Description Partially update subscription fields by ID (only provided fields are modified)
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID" minimum(1)
// @Param subscription body models.UpdateSubscriptionRequest true "Update payload (partial update)"
// @Success 200 {object} models.SubscriptionResponse
// @Failure 400 {object} models.ErrorResponse "Bad Request - Invalid input or validation failed"
// @Failure 404 {object} models.ErrorResponse "Not Found - Subscription does not exist"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {

	var reqUri *models.SubscriptionUriIDRequest

	// Bind and validate uri request payload
	//Привязка и проверка полезной нагрузки запроса URI
	if err := c.ShouldBindUri(&reqUri); err != nil {
		h.Logger.WithError(err).Info(validations.ErrInvalidRequestInput)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: validations.ErrInvalidRequestInput.Error(), Details: err.Error(),
		})
		return
	}
	// Bind and validate update request payload.
	// Привязать и проверить полезную нагрузку запроса на обновление.
	var req *models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.WithError(err).Info(validations.ErrInvalidRequestInput)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: validations.ErrInvalidRequestInput.Error(), Details: err.Error()})
		return
	}

	h.Logger.Info("updating subscription:")

	//process business logic for UpdateSubscriptionRequest
	//Обработка бизнес-логики для GetSubscription Request
	sub, err := h.service.UpdateSubscriptionByID(c.Request.Context(), reqUri.ID, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, FormatToSubscriptionResponse(sub))
}

// DeleteSubscription handles deleting a subscription by its ID.
// It validates the ID parameter, calls the repository to delete the record,
// logs any errors, and returns appropriate HTTP status codes.
// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Permanently delete a subscription by ID
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID" minimum(1)
// @Success 204 "No Content - Subscription successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Bad Request - Invalid subscription ID"
// @Failure 404 {object} models.ErrorResponse "Not Found - Subscription does not exist"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	var req *models.SubscriptionUriIDRequest
	// Bind and validate uri request payload
	//Привязка и проверка полезной нагрузки запроса URI
	if err := c.ShouldBindUri(&req); err != nil {
		h.Logger.WithError(err).Info(validations.ErrInvalidRequestInput)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: validations.ErrInvalidRequestInput.Error(), Details: err.Error(),
		})
		return
	}

	h.Logger.Info("deleting subscription by ID: ", req.ID)

	//process business logic for DeleteSubscriptionRequest
	//Обработка бизнес-логики для DeleteSubscription Request
	if err := h.service.DeleteSubscription(c.Request.Context(), req.ID); err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetUserSubscriptionSummary calculates subscription statistics for a given user
// within an optional date range and optional service name filter.
// Returns total cost, unique months, and subscription count.
// GetUserSubscriptionSummary godoc
// @Summary Get user subscription summary
// @Description Calculate subscription statistics including total cost, unique months, and count for a user
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param user_id query string true "User UUID" format(uuid)
// @Param service_name query string true "Filter by service name"
// @Param from query string false "Start date (MM-YYYY)"
// @Param to query string false "End date (MM-YYYY)"
// @Success 200 {object} models.UserSubscriptionSummaryResponse
// @Failure 400 {object} models.ErrorResponse "Bad Request - Invalid parameters"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /subscriptions/summary [get]
func (h *SubscriptionHandler) GetUserSubscriptionSummary(c *gin.Context) {

	var req *models.UserSubscriptionSummaryRequest

	// Bind and validate request payload
	//Привяжите и проверьте полезную нагрузку запроса.
	if err := c.ShouldBindQuery(&req); err != nil {
		h.Logger.WithError(err).Info(validations.ErrInvalidRequestInput)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: validations.ErrInvalidRequestInput.Error(), Details: err.Error()})
		return
	}

	h.Logger.Infof("getting user's subscription summary: UserID: %+v, ServiceName: %+v, PeriodStart: %+v, PeriodEnd: %+v", req.UserID, req.ServiceName, req.From, req.To)

	//process business logic for GetUserSubscriptionSummaryRequest
	//Обработка бизнес-логики для GetUserSubscriptionSummaryRequest
	unitPrice, totalAmount, totalMonths, err := h.service.GetUserSubscriptionSummary(c.Request.Context(), req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	res := &models.UserSubscriptionSummaryResponse{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		TotalMonths: totalMonths,
		UnitPrice:   unitPrice,
		TotalAmount: totalAmount,
	}
	c.JSON(http.StatusOK, res)

}

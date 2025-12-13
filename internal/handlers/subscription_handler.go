package handlers

import (
	"context"
	"net/http"

	"github.com/cyb3rkh4l1d/subsapi/internal/models"
	"github.com/cyb3rkh4l1d/subsapi/internal/repository"
	"github.com/cyb3rkh4l1d/subsapi/internal/utils"
	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SubscriptionHandler handles HTTP requests related to subscriptions.
// It contains shared context, logger, and repository dependencies.

type SubscriptionHandler struct {
	ctx           context.Context
	Logger        *logrus.Entry
	SubRepository *repository.SubscriptionRepository
}

// @Description Defines the request body for creating a new subscription.
type createSubReq struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,gt=0"`
	UserID      string  `json:"user_id" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

// @Description Defines the request body for updating a subscription.
type updateSubReq struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty" binding:"gt=0"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

// @Description Defines the API response structure for a subscription.
type SubscriptionResponse struct {
	ID          int    `json:"service_id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

//...................................................

// @Description Defines the generic error
type ErrorResponse struct {
	Error string `json:"error"`
}

// @Description Defines the API response structure for /stats endpoint.
type StatsResponse struct {
	Total int `json:"total"`
	Count int `json:"count"`
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/

// NewSubscriptionHandlers creates and returns a SubscriptionHandler instance with
// With shared context, logger, and repository dependencies.
func NewSubscriptionHandlers(ctx context.Context, handlerLogger *logrus.Entry, repo *repository.SubscriptionRepository) SubscriptionHandler {
	return SubscriptionHandler{ctx: ctx, Logger: handlerLogger, SubRepository: repo}
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
// @Param subscription body createSubReq true "Subscription payload"
// @Success 201 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {

	var req createSubReq

	// Bind and validate request payload
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.WithError(err).Warn("[-] invalid request payload in CreateSubscription")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inputs"})
		return
	}

	// Parse start_date (MM-YYYY)
	startDate, err := validations.ValidateStartDate(req.StartDate, h.Logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse optional end_date (MM-YYYY)
	endDate, err := validations.ValidateEndDate(startDate, *req.EndDate, h.Logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Validate Price
	if err := validations.ValidatePrice(req.Price, h.Logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Validate ServiceName
	if err := validations.ValidateServiceName(req.ServiceName, h.Logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build subscription model
	sub := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	// Persist subscription to database
	if err := h.SubRepository.CreateSubscription(h.ctx, sub); err != nil {
		errMsg := "Failed to create subscription"
		h.Logger.WithError(err).Error("[-] " + errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	// Convert model entity to API response format (MM-YYYY for dates) and
	// Return the normalized subscription list to the client
	res := ToResponse(sub)
	c.JSON(http.StatusCreated, res)
}

// ListSubscriptions returns all subscriptions in the system.
// It converts internal date fields to MM-YYYY format
// and responds with a normalized API payload.
// ListSubscriptions godoc
// @Summary List all subscriptions
// @Description Retrieve all subscriptions
// @Tags Subscriptions
// @Produce json
// @Success 200 {array} SubscriptionResponse
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	// Fetch all subscriptions from the repository
	subs, err := h.SubRepository.List(h.ctx)
	if err != nil {
		errMsg := "Error Getting Subscriptions List"
		h.Logger.WithError(err).Error("[-] " + errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	// Convert model entities to API response format (MM-YYYY for dates) and
	// Return the normalized subscription list to the client
	res := make([]SubscriptionResponse, len(subs))
	for i, sub := range subs {
		res[i] = ToResponse(&sub)
	}

	c.JSON(http.StatusOK, res)

}

// GetSubscription retrieves a single subscription by its ID.
// It validates the identifier and returns a formatted
// subscription response if found.
// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Retrieve a subscription using its ID
// @Tags Subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	// Extract and validate the subscription ID from the URL path
	id, err := validations.ValidateSubscriptionID(c.Param("id"), h.Logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the subscription by ID from the repository
	sub, err := h.SubRepository.GetByID(h.ctx, uint(id))
	if err != nil {
		errMsg := "Error Fetching Data By Id"
		h.Logger.WithError(err).Error("[-] " + errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	// Handle case where no subscription exists for the given ID
	if sub == nil {
		errMsg := "No Subscriptions For This ID"
		h.Logger.WithError(err).Error("[-] " + errMsg)
		c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		return
	}
	// Convert model entity to API response format and return it
	res := ToResponse(sub)
	c.JSON(http.StatusOK, res)

}

// UpdateSubscription updates an existing subscription by ID.
// Only fields provided in the request are modified,
// with validation applied to price and date formats.
// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update subscription fields by ID
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body updateSubReq true "Update payload"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	// Parse and validate subscription ID from URL path
	id, err := validations.ValidateSubscriptionID(c.Param("id"), h.Logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Bind and validate update request payload.
	var req updateSubReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := "Invalid Inputs"
		h.Logger.WithError(err).Error("[-] " + errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	// Fetch existing subscription to apply partial updates
	sub, err := h.SubRepository.GetByID(h.ctx, uint(id))
	if err != nil {
		errMsg := "Error Getting Subscription ID"
		h.Logger.WithError(err).Error("[-] " + errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	if sub == nil {
		errMsg := "Subscription Not Found"
		h.Logger.WithError(err).Error("[-] " + errMsg)
		c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		return
	}

	// Update service name if provided
	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}

	// Update price if provided and ensure it is non-negative
	if req.Price != nil {
		if err := validations.ValidatePrice(*req.Price, h.Logger); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sub.Price = *req.Price
	}
	// Update start date if provided and validate format
	if req.StartDate != nil {
		startDate, err := validations.ValidateStartDate(*req.StartDate, h.Logger)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sub.StartDate = startDate
	}
	// Update or clear end date and enforce end_date > start_date
	if req.EndDate != nil {
		if *req.EndDate == "" {
			// Explicitly clear end_date
			sub.EndDate = nil
		} else {
			endDate, err := validations.ValidateEndDate(sub.StartDate, *req.EndDate, h.Logger)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			sub.EndDate = endDate
		}
	}

	if err := h.SubRepository.Update(h.ctx, sub); err != nil {
		errMsg := "update failed"
		h.Logger.WithError(err).Error("[-] " + errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	// Return updated subscription in API response format
	res := ToResponse(sub)
	c.JSON(http.StatusOK, res)
}

// DeleteSubscription handles deleting a subscription by its ID.
// It validates the ID parameter, calls the repository to delete the record,
// logs any errors, and returns appropriate HTTP status codes:
// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags Subscriptions
// @Param id path int true "Subscription ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	// Validate the ID parameter; return 400 if not a positive integer
	id, err := validations.ValidateSubscriptionID(c.Param("id"), h.Logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call repository to delete the subscription; log and return 500 if an error occurs
	if err := h.SubRepository.Delete(h.ctx, uint(id)); err != nil {
		errMsg := "Failed To Delete The Subscription"
		h.Logger.WithError(err).Error("[-] " + errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	c.Status(http.StatusNoContent)

}

// SumCostHandler calculates the total cost of subscriptions for a given user
// within an optional date range and optional service name filter.
// It reads query parameters 'from' (start date,optional), 'to' (end date,optional), 'user_id' (required)
// and 'service_name' (optional), validates them, and returns the total cost in JSON.
// SumCostHandler godoc
// @Summary Calculate total subscription cost
// @Description Calculate total cost for a user within a period
// @Tags Subscriptions
// @Produce json
// @Param user_id query string true "User UUID"
// @Param from query string false "Start period (MM-YYYY)"
// @Param to query string false "End period (MM-YYYY)"
// @Param service_name query string false "Service name filter"
// @Success 200 {object} StatsResponse "Response for static"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /subscriptions/stats [get]
func (h *SubscriptionHandler) SumCostHandler(c *gin.Context) {
	// Read query parameters for filtering: from, to, subscription_name, and required user_id
	startStr := c.Query("from")
	endStr := c.Query("to")
	serviceName := c.Query("service_name")
	userID := c.Query("user_id")

	// Validate required user_id
	if err := validations.ValidateUserID(userID, h.Logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse 'from' date string to time.Time
	periodStart, err := validations.ValidateStartDateSumCostHandler(startStr, h.Logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse 'to' date string to time.Time
	periodEnd, err := validations.ValidateEndDateSumCostHandler(periodStart, endStr, h.Logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call repository to calculate the total cost for the period and filters
	total, count, err := h.SubRepository.CalculateTotalCost(h.ctx, periodStart, periodEnd, userID, serviceName)
	if err != nil {
		errMsg := "Failed To Calculate TotalCost"
		h.Logger.Error("[-] " + errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total, "count": count})

}

// ToResponse converts a Subscription model to a SubscriptionResponse DTO
// formatting the StartDate and EndDate in "MM-YYYY" format.
func ToResponse(sub *models.Subscription) SubscriptionResponse {
	// prepare EndDate string; leave empty if EndDate is nil or zero
	var end string
	if sub.EndDate != nil && !sub.EndDate.IsZero() {
		end = utils.FormatMonthYear(*sub.EndDate)
	}
	// return response object with formatted dates
	return SubscriptionResponse{
		ID:          int(sub.ID),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   utils.FormatMonthYear(sub.StartDate),
		EndDate:     end,
	}
}

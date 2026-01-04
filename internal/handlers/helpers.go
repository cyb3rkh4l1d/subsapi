package handlers

import (
	"net/http"

	"github.com/cyb3rkh4l1d/subsapi/internal/models"
	"github.com/cyb3rkh4l1d/subsapi/internal/utils"
	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/gin-gonic/gin"
)

// ToResponse converts a Subscription model to a SubscriptionResponse DTO
// formatting the StartDate and EndDate in "MM-YYYY" format.
// ToResponse преобразует модель Subscription в DTO SubscriptionResponse
// форматирование StartDate и EndDate в формате "MM-YYYY".
func FormatToSubscriptionResponse(sub *models.Subscription) models.SubscriptionResponse {
	var end string
	if sub.EndDate != nil && !sub.EndDate.IsZero() {
		end = utils.FormatMonthYear(*sub.EndDate)
	}
	// return response object with formatted dates
	// Возвращает объект ответа с отформатированными датами
	return models.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   utils.FormatMonthYear(sub.StartDate),
		EndDate:     end,
	}
}

// handleServiceError maps service layer errors to appropriate HTTP responses
// Функция handleServiceError сопоставляет ошибки уровня сервиса с соответствующими HTTP-ответами.
func (h *SubscriptionHandler) handleServiceError(c *gin.Context, err error) {
	switch err {
	case validations.ErrInvalidServiceName,
		validations.ErrEmptyUserID,
		validations.ErrInvalidPrice,
		validations.ErrInvalidDateFormat,
		validations.ErrInvalidStartDate,
		validations.ErrInvalidEndDate,
		validations.ErrEndDateBeforeStart,
		validations.ErrInvalidSubscriptionID,
		validations.ErrInvalidUserID:
		h.service.Logger.Info(err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	case validations.ErrSubscriptionNotFound:
		h.service.Logger.Info(err)
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
	case validations.ErrSubscriptionExists:
		h.service.Logger.Warn(err)
		c.JSON(http.StatusConflict, models.ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal server error"})
	}
}

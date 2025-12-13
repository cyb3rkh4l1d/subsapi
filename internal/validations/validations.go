package validations

import (
	"errors"
	"strconv"
	"time"

	"github.com/cyb3rkh4l1d/subsapi/internal/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ValidateUserID ensures the userID is not empty and is a valid UUID
func ValidateUserID(userID string, logger *logrus.Entry) error {
	if userID == "" {
		logger.Warn("user_id is empty")
		return errors.New("user_id must be provided")
	}
	if _, err := uuid.Parse(userID); err != nil {
		logger.WithError(err).Warn("invalid user_id UUID")
		return errors.New("invalid user_id UUID")
	}
	return nil
}

// ValidateServiceName ensures service name is not empty
func ValidateServiceName(name string, logger *logrus.Entry) error {
	if name == "" {
		logger.Warn("service name is empty")
		return errors.New("service_name must be provided")
	}
	return nil
}

// ValidatePrice ensures the price is positive
func ValidatePrice(price int, logger *logrus.Entry) error {
	if price <= 0 {
		logger.Warnf("invalid price: %d", price)
		return errors.New("price must be a positive integer")
	}
	return nil
}

// ValidateStartDate parses and validates a start date in MM-YYYY format
func ValidateStartDate(dateStr string, logger *logrus.Entry) (time.Time, error) {
	startDate, err := utils.ParseMonthYear(dateStr)
	if err != nil {
		logger.WithError(err).Warnf("invalid start_date: %s", dateStr)
		return time.Time{}, errors.New("invalid start_date format, expected MM-YYYY")
	}
	return startDate, nil
}

// ValidateEndDate parses and validates end date, ensures end > start if provided
func ValidateEndDate(startDate time.Time, endStr string, logger *logrus.Entry) (*time.Time, error) {
	if endStr == "" {
		return nil, nil
	}
	endDate, err := utils.ParseMonthYear(endStr)
	if err != nil {
		logger.WithError(err).Warnf("invalid end_date: %s", endStr)
		return nil, errors.New("invalid end_date format, expected MM-YYYY")
	}
	if !endDate.After(startDate) {
		logger.Warnf("end_date %v is not after start_date %v", endDate, startDate)
		return nil, errors.New("end_date must be after start_date")
	}
	return &endDate, nil
}

// ValidateSubscriptionID parses and validates subscription ID from URL path
func ValidateSubscriptionID(idStr string, logger *logrus.Entry) (uint, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logger.WithError(err).Warnf("invalid subscription ID: %s", idStr)
		return 0, errors.New("invalid subscription ID")
	}
	return uint(id), nil
}

// This is used in SumCostHandler
func ValidateStartDateSumCostHandler(dateStr string, logger *logrus.Entry) (time.Time, error) {
	startDate, err := utils.ParseMonthYearSumHandler(dateStr)
	if err != nil {
		logger.WithError(err).Warnf("invalid start_date: %s", dateStr)
		return time.Time{}, errors.New("invalid start_date format, expected MM-YYYY")
	}
	return startDate, nil
}

// ValidateEndDate parses and validates end date, ensures end > start if provided
func ValidateEndDateSumCostHandler(startDate time.Time, endStr string, logger *logrus.Entry) (time.Time, error) {
	if endStr == "" {
		return time.Time{}, nil
	}
	endDate, err := utils.ParseMonthYearSumHandler(endStr)
	if err != nil {
		logger.WithError(err).Warnf("invalid end_date: %s", endStr)
		return time.Time{}, errors.New("invalid end_date format, expected MM-YYYY")
	}
	if !endDate.After(startDate) {
		logger.Warnf("end_date %v is not after start_date %v", endDate, startDate)
		return time.Time{}, errors.New("end_date must be after start_date")
	}
	return endDate, nil
}

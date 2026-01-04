package validations

import (
	"time"

	"github.com/cyb3rkh4l1d/subsapi/internal/utils"
	"github.com/google/uuid"
)

// ValidateUserID ensures the userID is not empty and is a valid UUID
// Функция ValidateUserID гарантирует, что userID не пуст и является действительным UUID.
func ValidateUserID(userID string) error {
	if userID == "" {
		return ErrEmptyUserID
	}
	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}
	return nil
}

// ValidateServiceName ensures service name is not empty
// ValidateServiceName гарантирует, что имя сервиса не пустое
func ValidateServiceName(name string) error {
	if name == "" {
		return ErrInvalidServiceName
	}
	return nil
}

// ValidatePrice ensures the price is positive
// Функция ValidatePrice гарантирует, что цена положительная
func ValidatePrice(price int) error {
	if price <= 0 {
		return ErrInvalidPrice
	}
	return nil
}

// ValidateStartDate parses and validates a start date in MM-YYYY format
// Функция ValidateStartDate анализирует и проверяет дату начала в формате MM-YYYY.
func ValidateStartDate(dateStr string) (time.Time, error) {
	startDate, err := utils.ParseMonthYear(dateStr)
	if err != nil {
		return time.Time{}, ErrInvalidStartDate
	}
	return startDate, nil
}

// ValidateEndDate parses and validates end date, ensures end >= start if provided
// Функция ValidateEndate анализирует и проверяет дату окончания, обеспечивая, чтобы дата окончания была >= даты начала, если она указана.
func ValidateEndDate(startDate time.Time, endStr string) (*time.Time, error) {
	if endStr == "" {
		return nil, nil
	}
	endDate, err := utils.ParseMonthYear(endStr)
	if err != nil {
		return nil, ErrInvalidEndDate
	}

	if endDate.Before(startDate) {
		return nil, ErrEndDateBeforeStart
	}
	return &endDate, nil
}

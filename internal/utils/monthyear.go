package utils

import (
	"time"
)

// ParseMonthYear parses strings like "07-2025" into time.Time
// with day set to the first day of the month.
func ParseMonthYear(value string) (time.Time, error) {
	return time.Parse("01-2006", value)
}

// FormatMonthYear, convert time to mm-yyyy format
func FormatMonthYear(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("01-2006")
}

// ParseMonthYear parses strings like "07-2025" into time.Time
// Use only in SumCostHandler
func ParseMonthYearSumHandler(value string) (time.Time, error) {

	if value == "" {

		return time.Time{}, nil
	}
	return time.Parse("01-2006", value)
}

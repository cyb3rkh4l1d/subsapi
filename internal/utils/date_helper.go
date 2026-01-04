package utils

import (
	"time"
)

// ParseMonthYear parses strings like "07-2025" into time.Time
// with day set to the first day of the month.
// ParseMonthYear преобразует строки типа "07-2025" в time.Time
// где day устанавливается на первый день месяца.
func ParseMonthYear(value string) (time.Time, error) {
	return time.Parse("01-2006", value)
}

// FormatMonthYear, convert time to mm-yyyy format
// FormatMonthYear, преобразование времени в формат мм-гггг
func FormatMonthYear(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("01-2006")
}

// MaxTime returns the later of two time values (maximum)
// Функция MaxTime возвращает более позднее из двух значений времени (максимальное значение).
func MaxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// MinTime returns the earlier of two time values (minimum)
// MinTime возвращает более раннее из двух значений времени (минимум).
func MinTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

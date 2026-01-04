package service

import (
	"fmt"
	"time"

	"github.com/cyb3rkh4l1d/subsapi/internal/models"
	"github.com/cyb3rkh4l1d/subsapi/internal/utils"
)

// Calculates total cost: Sum of (monthly price × months active within period)
// Counts unique months: Deduplicates months when multiple subscriptions overlap
// Вычисляет общую стоимость: Сумма (месячная цена × количество активных месяцев в течение периода)
// Подсчитывает уникальные месяцы: Удаляет дубликаты месяцев, если несколько подписок перекрываются
func CalculateSubscriptionMetrics(
	subscriptions []models.Subscription,
	periodStart time.Time, periodEnd time.Time,
) (int, int64, int) {
	var unitPrice int
	var totalCost int64
	uniqueMonths := make(map[string]bool)

	// Process each subscription the user has
	// Обработка каждой подписки, имеющейся у пользователя
	for _, sub := range subscriptions {
		//set unitPrice once
		if unitPrice == 0 {
			unitPrice = sub.Price
		}
		var effectiveEnd time.Time
		//return largest date between subscription startDate and query from/periodStart
		// e.g if subscription starts Mar 2024, but query starts Jan 2024, overlap starts at Mar 2024
		//Возвращает наибольшую дату между датой начала подписки и датой начала запроса/периодом начала запроса.
		//Например, если подписка начинается в марте 2024 года, а запрос — в январе 2024 года, то совпадение начинается с марта 2024 года.
		effectiveStart := utils.MaxTime(sub.StartDate, periodStart)
		if sub.EndDate == nil || sub.EndDate.IsZero() {
			effectiveEnd = periodEnd // Ongoing subscription
		} else {
			//return least date between subscription startDate and query from/periodStart
			//If subscription ends Jul 2024, but query ends Jun 2024, overlap ends at Jun 2024 (the earlier date).
			//Возвращает наименьшую дату между датой начала подписки и датой начала запроса/периодом начала
			//Если подписка заканчивается в июле 2024 года, а запрос — в июне 2024 года, то перекрытие заканчивается в июне 2024 года (более ранняя дата).
			effectiveEnd = utils.MinTime(*sub.EndDate, periodEnd)

		}

		// Check if there's any overlap
		// Проверьте, нет ли пересечений
		if effectiveStart.After(effectiveEnd) || effectiveStart.Equal(effectiveEnd.AddDate(0, 0, 1)) {
			continue // No overlap
		}

		// Calculate months and add to unique set
		// Рассчитать количество месяцев и добавить их в уникальный набор
		monthsAdded := AddOverlapMonths(uniqueMonths, effectiveStart, effectiveEnd)

		if monthsAdded > 0 {
			// Calculate cost for these months
			// Рассчитать стоимость за эти месяцы
			subscriptionCost := int64(sub.Price) * int64(monthsAdded)
			totalCost += subscriptionCost
		}
	}

	return unitPrice, totalCost, len(uniqueMonths)
}

// Calculates how many months between effectiveStart and effectiveEnd
// Adds each month to the uniqueMonths map (deduplicates automatically)
// Вычисляет количество месяцев между effectiveStart и effectiveEnd
// Добавляет каждый месяц в карту uniqueMonths (автоматически удаляет дубликаты)
func AddOverlapMonths(
	uniqueMonths map[string]bool,
	start, end time.Time,
) int {

	current := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
	endMonth := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location())

	monthsAdded := 0

	// Iterate through each month from the range current-endMonth
	//update the map if key does'nt exist in the map
	// Проходим по каждому месяцу в диапазоне current-endMonth
	// // Обновляем карту, если ключ отсутствует в карте
	for !current.After(endMonth) {
		monthKey := fmt.Sprintf("%d-%02d", current.Year(), current.Month())
		if !uniqueMonths[monthKey] {
			uniqueMonths[monthKey] = true
			monthsAdded++
		}
		current = current.AddDate(0, 1, 0) // Next month. В следующем месяце
	}

	return monthsAdded
}

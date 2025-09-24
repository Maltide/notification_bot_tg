package problems

import "time"

func FindSmallestTime(slice []time.Time) time.Time {
	// Эта функция находит самое раннее время в слайсе временных меток.
	//
	// Аргумент:
	//   - slice []time.Time: слайс временных меток
	//   - например: []time.Time{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)}
	//
	// Алгоритм:
	//   - Придумай самостоятельно, используя логику похожую на поиск минимального числа
	//
	// Что может понадобиться:
	//   - Для сравнения двух временных меток t1 и t2 используй метод t1.Before(t2)
	//   - t1.Before(t2) возвращает true, если t1 раньше чем t2
	//
	// Пример для слайса с временами 12:00 и 10:00:
	//   - Возвращаем 10:00, так как это самое раннее время в слайсе
	if len(slice) == 0 {
		return time.Time{}
	}
	minTime := slice[0]

	for i := range slice[1:] {
		if slice[i].Before(minTime) {
			minTime = slice[i]
		}
	}

	return minTime
}

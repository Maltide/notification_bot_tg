package problems

func FindSmallestInt(slice []int) int {
	// Эта функция находит наименьшее число в слайсе целых чисел.
	//
	// Аргумент:
	//   - slice []int: слайс целых чисел, например [5, 3, 8, 1, 9]
	//
	// Алгоритм:
	//   1. Заводим переменную smallest, которая будет хранить наименьшее значение
	//   2. Устанавливаем smallest равным первому элементу слайса (slice[0])
	//   3. Проходимся по всем остальным элементам слайса начиная с индекса 1
	//   4. Сравниваем каждый элемент с smallest
	//   5. Если текущий элемент меньше smallest, обновляем smallest
	//   6. Возвращаем smallest
	//
	// Пример для слайса [5, 3, 8, 1, 9]:
	//   - smallest = 5 (slice[0])
	//   - slice[1] = 3, 3 < 5, обновляем smallest = 3
	//   - slice[2] = 8, 8 > 3, не обновляем
	//   - slice[3] = 1, 1 < 3, обновляем smallest = 1
	//   - slice[4] = 9, 9 > 1, не обновляем
	//   - Возвращаем 1 (это slice[3])
	if len(slice) == 0 {
		return 0
	}
	var smallest int

	smallest = slice[0]

	for i := range slice {

		if slice[i] < smallest {
			smallest = slice[i]
		}
	}

	return smallest
}

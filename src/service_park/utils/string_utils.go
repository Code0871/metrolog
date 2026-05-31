// ============================================================================
// Пакет utils предоставляет вспомогательные функции
//
// Назначение:
//   - ParseCommaSeparatedString - разбор строки с разделителями
//   - TrimStrings - очистка массива строк от пробелов
//
// ============================================================================
package utils

import "strings"

// ParseCommaSeparatedString - разбирает строку вида "1,2,3" в массив строк
// Очищает пробелы и удаляет пустые значения
//
// Пример:
//   ParseCommaSeparatedString("14780, 14781, 14783") // возвращает ["14780", "14781", "14783"]
//   ParseCommaSeparatedString("")                     // возвращает []string{}
func ParseCommaSeparatedString(input string) []string {
	if input == "" {
		return []string{}
	}

	rawItems := strings.Split(input, ",")
	result := make([]string, 0, len(rawItems))

	for _, item := range rawItems {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

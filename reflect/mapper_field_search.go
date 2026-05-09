package reflect

import (
	"strings"
)

func DefaultStringSearch(fieldName string) []string {
	result := make([]string, 5)
	result[0] = fieldName
	result[1] = strings.ToLower(fieldName)
	result[2] = strings.ToUpper(fieldName)
	result[3] = ToLowerSnakeCase(fieldName)
	result[4] = ToUpperSnakeCase(fieldName)
	return result
}

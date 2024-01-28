package reflect

import (
	"log/slog"
	"reflect"
	"strconv"
)

func GetString(raw any) string {
	input, isValid := GetValue(raw)
	if !isValid {
		return ""
	}

	var stringValue string
	switch input := input.(type) {
	case string:
		stringValue = input
	case int, int8, int16, int32, int64:
		inputValue := GetIntBase(input, 64)
		stringValue = strconv.FormatInt(inputValue, 10)
	case uint, uint8, uint16, uint32, uint64:
		uintValue := GetUintBase(input, 64)
		stringValue = strconv.FormatUint(uintValue, 10)
	case bool:
		stringValue = strconv.FormatBool(input)
	case float32, float64:
		floatValue := GetFloatBase(input, 64)
		stringValue = strconv.FormatFloat(floatValue, 'f', -1, 64)
	default:
		slog.Warn("conversion for string type failed", "type", reflect.TypeOf(input), "value", input)
	}

	return stringValue
}

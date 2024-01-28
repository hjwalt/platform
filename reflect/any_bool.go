package reflect

import (
	"log/slog"
	"reflect"
	"strconv"
	"strings"
)

func GetBool(raw any) bool {
	input, isValid := GetValue(raw)
	if !isValid {
		return false
	}
	var boolValue = false
	switch input := input.(type) {
	case bool:
		boolValue = input
	case string:
		var err error
		if len(input) == 0 {
			return false
		}
		boolValue, err = strconv.ParseBool(strings.ToUpper(input))
		if err != nil {
			slog.Warn("string parse bool failed", "error", err)
		}
	case int, int8, int16, int32, int64:
		intValue := GetIntBase(input, 64)
		switch intValue {
		case 0:
			boolValue = false
		case 1:
			boolValue = true
		default:
			boolValue = false
		}
	case uint, uint8, uint16, uint32, uint64:
		uintValue := GetUintBase(input, 64)
		switch uintValue {
		case 0:
			boolValue = false
		case 1:
			boolValue = true
		default:
			boolValue = false
		}
	default:
		slog.Warn("conversion for bool type failed", "type", reflect.TypeOf(input), "value", input)
	}
	return boolValue
}

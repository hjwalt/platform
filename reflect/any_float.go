package reflect

import (
	"log/slog"
	"reflect"
	"strconv"
)

func GetFloatBase(raw any, bitSize int) float64 {
	input, isValid := GetValue(raw)
	if !isValid {
		return float64(0)
	}
	var floatValue float64
	switch input := input.(type) {
	case float32:
		floatValue = float64(input)
	case float64:
		floatValue = input
	case int, int8, int16, int32, int64:
		inputValue := GetIntBase(input, 64)
		floatValue = float64(inputValue)
	case uint, uint8, uint16, uint32, uint64:
		uintValue := GetUintBase(input, 64)
		floatValue = float64(uintValue)
	case string:
		if input == "" {
			input = "0"
		}
		var err error
		floatValue, err = strconv.ParseFloat(input, bitSize)
		if err != nil {
			slog.Warn("string parse float failed", "error", err)
			return 0
		}
	default:
		slog.Warn("conversion for float type failed", "type", reflect.TypeOf(input), "value", input)
	}

	return floatValue
}

func GetFloat32(input any) float32 {
	return float32(GetFloatBase(input, 32))
}

func GetFloat64(input any) float64 {
	return float64(GetFloatBase(input, 64))
}

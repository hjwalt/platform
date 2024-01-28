package reflect

import (
	"log/slog"
	"math"
	"reflect"
	"strconv"
)

func GetIntBase(raw any, bitSize int) int64 {
	input, isValid := GetValue(raw)
	if !isValid {
		return int64(0)
	}
	var intValue int64
	switch input := input.(type) {
	case int:
		intValue = int64(input)
	case int8:
		intValue = int64(input)
	case int16:
		intValue = int64(input)
	case int32:
		intValue = int64(input)
	case int64:
		intValue = input
	case uint, uint8, uint16, uint32, uint64:
		uintValue := GetUintBase(input, 64)
		intValue = int64(uintValue)
	case float32, float64:
		intValue = int64(math.Round(GetFloatBase(input, 64)))
	case bool:
		if input {
			intValue = 1
		} else {
			intValue = 0
		}
	case string:
		if input == "" {
			input = "0"
		}
		var err error
		intValue, err = strconv.ParseInt(input, 10, bitSize)
		if err != nil {
			slog.Warn("string parse int failed", "error", err)
			return 0
		}
	case []byte:
		if len(input) == 0 {
			return 0
		}
		return int64(Endian().Uint64(input))
	default:
		slog.Warn("conversion for int type failed", "type", reflect.TypeOf(input), "value", input)
		return 0
	}

	return intValue
}

func GetInt(input any) int {
	return int(GetIntBase(input, 0))
}

func GetInt8(input any) int8 {
	return int8(GetIntBase(input, 8))
}

func GetInt16(input any) int16 {
	return int16(GetIntBase(input, 16))
}

func GetInt32(input any) int32 {
	return int32(GetIntBase(input, 32))
}

func GetInt64(input any) int64 {
	return int64(GetIntBase(input, 64))
}

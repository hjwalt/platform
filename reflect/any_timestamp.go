package reflect

import (
	"log/slog"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var UTC *time.Location

func init() {
	UTC, _ = time.LoadLocation("UTC")
}

func GetTimestampMs(raw any) *time.Time {
	input, isValid := GetValue(raw)
	if !isValid {
		return nil
	}

	var timeValue time.Time
	switch input := input.(type) {
	case time.Time:
		timeValue = input
	case timestamppb.Timestamp:
		timeValue = input.AsTime().UTC()
	case int, int8, int16, int32, int64:
		intValue := GetInt64(input)
		sec := intValue / 1000
		msec := intValue % 1000
		timeValue = time.Unix(sec, msec*int64(time.Millisecond)).UTC()
	case string:
		if timeWithNano, err := time.Parse(time.RFC3339Nano, input); err == nil {
			timeValue = timeWithNano.In(UTC)
		} else {
			slog.Warn("conversion parsing for timestamp type failed", "type", reflect.TypeOf(input), "value", input)
			return nil
		}
	default:
		slog.Warn("conversion for time type failed", "type", reflect.TypeOf(input), "value", input)
		return nil
	}

	return &timeValue
}

func GetProtoTimestampMs(raw any) *timestamppb.Timestamp {
	parsedTime := GetTimestampMs(raw)
	if parsedTime != nil {
		return timestamppb.New(*parsedTime)
	}
	return nil
}

package brave_search_web

import (
	"fmt"
	"time"

	brave "github.com/hjwalt/platform/agent/tool/brave_search"
)

func WithSubscriptionToken(header string) brave.Header {
	return brave.Header{
		Key:   "X-Subscription-Token",
		Value: header,
	}
}

func WithNoCache() brave.Header {
	return brave.Header{
		Key:   "Cache-Control",
		Value: "no-cache",
	}
}

func WithUserAgent(header string) brave.Header {
	return brave.Header{
		Key:   "User-Agent",
		Value: header,
	}
}

func WithLocLatitude(header float32) brave.Header {
	return brave.Header{
		Key:   "X-Loc-Lat",
		Value: fmt.Sprintf("%.3f", header),
	}
}

func WithLocLongitude(header float32) brave.Header {
	return brave.Header{
		Key:   "X-Loc-Long",
		Value: fmt.Sprintf("%.3f", header),
	}
}

func WithLocTimezone(header time.Location) brave.Header {
	return brave.Header{
		Key:   "X-Loc-Timezone",
		Value: header.String(),
	}
}

func WithLocCity(header string) brave.Header {
	return brave.Header{
		Key:   "X-Loc-City",
		Value: header,
	}
}

func WithLocState(header string) brave.Header {
	return brave.Header{
		Key:   "X-Loc-State",
		Value: header,
	}
}

func WithLocStateName(header string) brave.Header {
	return brave.Header{
		Key:   "X-Loc-State-Name",
		Value: header,
	}
}

func WithLocCountry(header string) brave.Header {
	return brave.Header{
		Key:   "X-Loc-Country",
		Value: header,
	}
}

func WithLocPostalCode(header string) brave.Header {
	return brave.Header{
		Key:   "X-Loc-Postal-Code",
		Value: header,
	}
}

func WithApiVersion(header string) brave.Header {
	return brave.Header{
		Key:   "Api-Version",
		Value: header,
	}
}

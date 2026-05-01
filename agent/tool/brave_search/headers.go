package brave_search

import (
	"fmt"
	"time"
)

func WithSubscriptionToken(header string) Header {
	return Header{
		Key:   "X-Subscription-Token",
		Value: header,
	}
}

func WithNoCache() Header {
	return Header{
		Key:   "Cache-Control",
		Value: "no-cache",
	}
}

func WithUserAgent(header string) Header {
	return Header{
		Key:   "User-Agent",
		Value: header,
	}
}

func WithLocLatitude(header float32) Header {
	return Header{
		Key:   "X-Loc-Lat",
		Value: fmt.Sprintf("%.3f", header),
	}
}

func WithLocLongitude(header float32) Header {
	return Header{
		Key:   "X-Loc-Long",
		Value: fmt.Sprintf("%.3f", header),
	}
}

func WithLocTimezone(header time.Location) Header {
	return Header{
		Key:   "X-Loc-Timezone",
		Value: header.String(),
	}
}

func WithLocCity(header string) Header {
	return Header{
		Key:   "X-Loc-City",
		Value: header,
	}
}

func WithLocState(header string) Header {
	return Header{
		Key:   "X-Loc-State",
		Value: header,
	}
}

func WithLocStateName(header string) Header {
	return Header{
		Key:   "X-Loc-State-Name",
		Value: header,
	}
}

func WithLocCountry(header string) Header {
	return Header{
		Key:   "X-Loc-Country",
		Value: header,
	}
}

func WithLocPostalCode(header string) Header {
	return Header{
		Key:   "X-Loc-Postal-Code",
		Value: header,
	}
}

func WithApiVersion(header string) Header {
	return Header{
		Key:   "Api-Version",
		Value: header,
	}
}

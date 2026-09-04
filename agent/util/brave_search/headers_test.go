package brave_search_test

import (
	"testing"
	"time"

	"github.com/hjwalt/platform/agent/util/brave_search"
	"github.com/stretchr/testify/assert"
)

func assertHeader(t *testing.T, header brave_search.Header, key, value string) {
	t.Helper()
	assert.Equal(t, key, header.Key)
	assert.Equal(t, value, header.Value)
}

func TestWithSubscriptionToken(t *testing.T) {
	assertHeader(t, brave_search.WithSubscriptionToken("tok-123"), "X-Subscription-Token", "tok-123")
}

func TestWithNoCache(t *testing.T) {
	assertHeader(t, brave_search.WithNoCache(), "Cache-Control", "no-cache")
}

func TestWithUserAgent(t *testing.T) {
	assertHeader(t, brave_search.WithUserAgent("curl/8.0"), "User-Agent", "curl/8.0")
}

func TestWithLocLatitude(t *testing.T) {
	assertHeader(t, brave_search.WithLocLatitude(float32(37.75)), "X-Loc-Lat", "37.750")
	assertHeader(t, brave_search.WithLocLatitude(float32(-12.5)), "X-Loc-Lat", "-12.500")
}

func TestWithLocLongitude(t *testing.T) {
	assertHeader(t, brave_search.WithLocLongitude(float32(-122.25)), "X-Loc-Long", "-122.250")
	assertHeader(t, brave_search.WithLocLongitude(float32(0)), "X-Loc-Long", "0.000")
}

func TestWithLocTimezone(t *testing.T) {
	loc := time.FixedZone("Custom/Zone", 2*60*60)
	assertHeader(t, brave_search.WithLocTimezone(*loc), "X-Loc-Timezone", "Custom/Zone")
}

func TestWithLocCity(t *testing.T) {
	assertHeader(t, brave_search.WithLocCity("Berkeley"), "X-Loc-City", "Berkeley")
}

func TestWithLocState(t *testing.T) {
	assertHeader(t, brave_search.WithLocState("CA"), "X-Loc-State", "CA")
}

func TestWithLocStateName(t *testing.T) {
	assertHeader(t, brave_search.WithLocStateName("California"), "X-Loc-State-Name", "California")
}

func TestWithLocCountry(t *testing.T) {
	assertHeader(t, brave_search.WithLocCountry("US"), "X-Loc-Country", "US")
}

func TestWithLocPostalCode(t *testing.T) {
	assertHeader(t, brave_search.WithLocPostalCode("94704"), "X-Loc-Postal-Code", "94704")
}

func TestWithApiVersion(t *testing.T) {
	assertHeader(t, brave_search.WithApiVersion("2024-04-15"), "Api-Version", "2024-04-15")
}

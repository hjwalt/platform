package brave_search_test

import (
	"testing"

	"github.com/hjwalt/platform/agent/util/brave_search"
	"github.com/stretchr/testify/assert"
)

func assertParam(t *testing.T, param brave_search.Param, key, value string) {
	t.Helper()
	assert.Equal(t, key, param.Key)
	assert.Equal(t, value, param.Value)
}

func TestWithCount(t *testing.T) {
	assertParam(t, brave_search.WithCount(10), "count", "10")
	assertParam(t, brave_search.WithCount(0), "count", "0")
}

func TestWithCountry(t *testing.T) {
	assertParam(t, brave_search.WithCountry("us"), "country", "us")
}

func TestWithExtraSnippets(t *testing.T) {
	assertParam(t, brave_search.WithExtraSnippets(true), "extra_snippets", "true")
	assertParam(t, brave_search.WithExtraSnippets(false), "extra_snippets", "false")
}

func TestWithFreshness(t *testing.T) {
	assertParam(t, brave_search.WithFreshness("py"), "freshness", "py")
}

func TestWithGogglesID(t *testing.T) {
	assertParam(t, brave_search.WithGogglesID("my-goggles"), "goggles_id", "my-goggles")
}

func TestWithOffset(t *testing.T) {
	assertParam(t, brave_search.WithOffset(5), "offset", "5")
}

func TestWithResultFilter(t *testing.T) {
	assertParam(t, brave_search.WithResultFilter([]string{"news", "videos"}), "result_filter", "news,videos")
	assertParam(t, brave_search.WithResultFilter(nil), "result_filter", "")
	assertParam(t, brave_search.WithResultFilter([]string{"single"}), "result_filter", "single")
}

func TestWithSafesearch(t *testing.T) {
	assertParam(t, brave_search.WithSafesearch("moderate"), "safesearch", "moderate")
}

func TestWithSearchLang(t *testing.T) {
	assertParam(t, brave_search.WithSearchLang("en"), "search_lang", "en")
}

func TestWithSpellcheck(t *testing.T) {
	assertParam(t, brave_search.WithSpellcheck(true), "spellcheck", "true")
	assertParam(t, brave_search.WithSpellcheck(false), "spellcheck", "false")
}

func TestWithTerm(t *testing.T) {
	assertParam(t, brave_search.WithTerm("golang"), "q", "golang")
}

func TestWithTextDecorations(t *testing.T) {
	assertParam(t, brave_search.WithTextDecorations(true), "text_decorations", "true")
	assertParam(t, brave_search.WithTextDecorations(false), "text_decorations", "false")
}

func TestWithUILang(t *testing.T) {
	assertParam(t, brave_search.WithUILang("en-US"), "ui_lang", "en-US")
}

func TestWithUnits(t *testing.T) {
	assertParam(t, brave_search.WithUnits("metric"), "units", "metric")
}

func TestWithSummary(t *testing.T) {
	assertParam(t, brave_search.WithSummary(true), "summary", "true")
	assertParam(t, brave_search.WithSummary(false), "summary", "false")
}

func TestWithSpellCheckLang(t *testing.T) {
	assertParam(t, brave_search.WithSpellCheckLang("fr"), "lang", "fr")
}

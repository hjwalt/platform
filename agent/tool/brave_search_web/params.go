package brave_search_web

import (
	"strings"

	brave "github.com/hjwalt/platform/agent/tool/brave_search"
	"github.com/hjwalt/platform/reflect"
)

func WithCount(param int) brave.Param {
	return brave.Param{
		Key:   "count",
		Value: reflect.GetString(param),
	}
}

func WithCountry(param string) brave.Param {
	return brave.Param{
		Key:   "country",
		Value: param,
	}
}

func WithExtraSnippets(param bool) brave.Param {
	return brave.Param{
		Key:   "extra_snippets",
		Value: reflect.GetString(param),
	}
}

func WithFreshness(param string) brave.Param {
	return brave.Param{
		Key:   "freshness",
		Value: param,
	}
}

func WithGogglesID(param string) brave.Param {
	return brave.Param{
		Key:   "goggles_id",
		Value: param,
	}
}

func WithOffset(param int) brave.Param {
	return brave.Param{
		Key:   "offset",
		Value: reflect.GetString(param),
	}
}

func WithResultFilter(param []string) brave.Param {
	return brave.Param{
		Key:   "result_filter",
		Value: strings.Join(param, ","),
	}
}

func WithSafesearch(param string) brave.Param {
	return brave.Param{
		Key:   "safesearch",
		Value: param,
	}
}

func WithSearchLang(param string) brave.Param {
	return brave.Param{
		Key:   "search_lang",
		Value: param,
	}
}

func WithSpellcheck(param bool) brave.Param {
	return brave.Param{
		Key:   "spellcheck",
		Value: reflect.GetString(param),
	}
}

func WithTerm(param string) brave.Param {
	return brave.Param{
		Key:   "q",
		Value: param,
	}
}

func WithTextDecorations(param bool) brave.Param {
	return brave.Param{
		Key:   "text_decorations",
		Value: reflect.GetString(param),
	}
}

func WithUILang(param string) brave.Param {
	return brave.Param{
		Key:   "ui_lang",
		Value: param,
	}
}

func WithUnits(param string) brave.Param {
	return brave.Param{
		Key:   "units",
		Value: param,
	}
}

func WithSummary(param bool) brave.Param {
	return brave.Param{
		Key:   "summary",
		Value: reflect.GetString(param),
	}
}

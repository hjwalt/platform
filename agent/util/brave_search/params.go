package brave_search

import (
	"strings"

	"github.com/hjwalt/platform/reflect"
)

func WithCount(param int) Param {
	return Param{
		Key:   "count",
		Value: reflect.GetString(param),
	}
}

func WithCountry(param string) Param {
	return Param{
		Key:   "country",
		Value: param,
	}
}

func WithExtraSnippets(param bool) Param {
	return Param{
		Key:   "extra_snippets",
		Value: reflect.GetString(param),
	}
}

func WithFreshness(param string) Param {
	return Param{
		Key:   "freshness",
		Value: param,
	}
}

func WithGogglesID(param string) Param {
	return Param{
		Key:   "goggles_id",
		Value: param,
	}
}

func WithOffset(param int) Param {
	return Param{
		Key:   "offset",
		Value: reflect.GetString(param),
	}
}

func WithResultFilter(param []string) Param {
	return Param{
		Key:   "result_filter",
		Value: strings.Join(param, ","),
	}
}

func WithSafesearch(param string) Param {
	return Param{
		Key:   "safesearch",
		Value: param,
	}
}

func WithSearchLang(param string) Param {
	return Param{
		Key:   "search_lang",
		Value: param,
	}
}

func WithSpellcheck(param bool) Param {
	return Param{
		Key:   "spellcheck",
		Value: reflect.GetString(param),
	}
}

func WithTerm(param string) Param {
	return Param{
		Key:   "q",
		Value: param,
	}
}

func WithTextDecorations(param bool) Param {
	return Param{
		Key:   "text_decorations",
		Value: reflect.GetString(param),
	}
}

func WithUILang(param string) Param {
	return Param{
		Key:   "ui_lang",
		Value: param,
	}
}

func WithUnits(param string) Param {
	return Param{
		Key:   "units",
		Value: param,
	}
}

func WithSummary(param bool) Param {
	return Param{
		Key:   "summary",
		Value: reflect.GetString(param),
	}
}

func WithSpellCheckLang(param string) Param {
	return Param{
		Key:   "lang",
		Value: param,
	}
}

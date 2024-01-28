package htmx

import (
	"net/http"
	"net/url"

	"github.com/hjwalt/platform/reflect"
)

var (
	HXBoosted               string = "HX-Boosted"
	HXCurrentUrl            string = "HX-Current-Url"
	HXHistoryRestoreRequest string = "HX-History-Restore-Request"
	HXPrompt                string = "HX-Prompt"
	HXRequest               string = "HX-Request"
	HXTarget                string = "HX-Target"
	HXTriggerName           string = "HX-Trigger-Name"
)

type HxRequestHeader struct {
	HxBoosted               bool
	HxRootURL               string
	HxCurrentURL            string
	HxHistoryRestoreRequest bool
	HxPrompt                string
	HxRequest               bool
	HxTarget                string
	HxTriggerName           string
	HxTrigger               string
}

func Extract(r *http.Request) HxRequestHeader {
	currentUrl := r.Header.Get(HXCurrentUrl)
	currParsed, _ := url.Parse(currentUrl)

	return HxRequestHeader{
		HxBoosted:               reflect.GetBool(r.Header.Get(HXBoosted)),
		HxCurrentURL:            currentUrl,
		HxRootURL:               currParsed.Scheme + "://" + currParsed.Host,
		HxHistoryRestoreRequest: reflect.GetBool(r.Header.Get(HXHistoryRestoreRequest)),
		HxPrompt:                r.Header.Get(HXPrompt),
		HxRequest:               reflect.GetBool(r.Header.Get(HXRequest)),
		HxTarget:                r.Header.Get(HXTarget),
		HxTriggerName:           r.Header.Get(HXTriggerName),
		HxTrigger:               r.Header.Get(HXTrigger),
	}
}

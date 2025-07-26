package htmx

import (
	"encoding/json"
	"net/http"

	"github.com/hjwalt/platform/commons/reflect"
)

var (
	HXLocation           string = "HX-Location"             // Allows you to do a client-side redirect that does not do a full page reload
	HXPushUrl            string = "HX-Push-Url"             // pushes a new url into the history stack
	HXRedirect           string = "HX-Redirect"             // can be used to do a client-side redirect to a new location
	HXRefresh            string = "HX-Refresh"              // if set to "true" the client side will do a full refresh of the page
	HXReplaceUrl         string = "HX-Replace-Url"          // replaces the current URL in the location bar
	HXReswap             string = "HX-Reswap"               // Allows you to specify how the response will be swapped. See hx-swap for possible values
	HXRetarget           string = "HX-Retarget"             // A CSS selector that updates the target of the content update to a different element on the page
	HXTrigger            string = "HX-Trigger"              // allows you to trigger client side events, see the documentation for more info
	HXTriggerAfterSettle string = "HX-Trigger-After-Settle" // allows you to trigger client side events, see the documentation for more info
	HXTriggerAfterSwap   string = "HX-Trigger-After-Swap"   // allows you to trigger client side events, see the documentation for more info
)

type LocationInput struct {
	Source  string                 `json:"source"`  // source - the source element of the request
	Event   string                 `json:"event"`   //event - an event that "triggered" the request
	Handler string                 `json:"handler"` //handler - a callback that will handle the response HTML
	Target  string                 `json:"target"`  //target - the target to swap the response into
	Swap    string                 `json:"swap"`    //swap - how the response will be swapped in relative to the target
	Values  map[string]interface{} `json:"values"`  //values - values to submit with the request
	Header  map[string]interface{} `json:"headers"` //headers - headers to submit with the request

}

// Location can be used to trigger a client side redirection without reloading the whole page
// https://htmx.org/headers/hx-location/
func Location(w http.ResponseWriter, li *LocationInput) error {
	payload, err := json.Marshal(li)
	if err != nil {
		return err
	}

	w.Header().Set(HXLocation, string(payload))
	return nil
}

// PushURL pushes a new url into the history stack.
// https://htmx.org/headers/hx-push-url/
func PushURL(w http.ResponseWriter, val string) {
	w.Header().Set(HXPushUrl, val)
}

// Redirect can be used to do a client-side redirect to a new location
func Redirect(w http.ResponseWriter, val string) {
	w.Header().Set(HXRedirect, val)
}

// Refresh if set to true the client side will do a full refresh of the page
func Refresh(w http.ResponseWriter, val bool) {
	w.Header().Set(HXRefresh, reflect.GetString(val))
}

// ReplaceURL allows you to replace the current URL in the browser location history.
// https://htmx.org/headers/hx-replace-url/
func ReplaceURL(w http.ResponseWriter, val string) {
	w.Header().Set(HXReplaceUrl, val)
}

// ReSwap allows you to specify how the response will be swapped. See hx-swap for possible values
// https://htmx.org/attributes/hx-swap/
func ReSwap(w http.ResponseWriter, val string) {
	w.Header().Set(HXReswap, val)
}

// ReTarget a CSS selector that updates the target of the content update to a different element on the page
func ReTarget(w http.ResponseWriter, val string) {
	w.Header().Set(HXRetarget, val)
}

// Trigger triggers events as soon as the response is received.
// https://htmx.org/headers/hx-trigger/
func Trigger(w http.ResponseWriter, val string) {
	w.Header().Set(HXTrigger, val)
}

// TriggerAfterSettle trigger events after the settling step.
// https://htmx.org/headers/hx-trigger/
func TriggerAfterSettle(w http.ResponseWriter, val string) {
	w.Header().Set(HXTriggerAfterSettle, val)
}

// TriggerAfterSwap trigger events after the swap step.
// https://htmx.org/headers/hx-trigger/
func TriggerAfterSwap(w http.ResponseWriter, val string) {
	w.Header().Set(HXTriggerAfterSwap, val)
}

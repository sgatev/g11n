package http

import (
	"net/http"

	"github.com/s2gatev/g11n"

	"golang.org/x/text/language"
)

// SetLocale sets the locale of a MessageFactory from HTTP Request value.
func SetLocale(mf *g11n.MessageFactory, r *http.Request) {
	acceptLanguage := r.Header.Get("Accept-Language")
	preferred, _, _ := language.ParseAcceptLanguage(acceptLanguage)

	var matcher = language.NewMatcher(mf.Locales())
	tag, _, _ := matcher.Match(preferred...)

	mf.LoadLocale(tag)
}

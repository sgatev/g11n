package http_test

import (
	"net/http"
	"testing"

	"golang.org/x/text/language"

	. "github.com/s2gatev/g11n"
	. "github.com/s2gatev/g11n/http"
	. "github.com/s2gatev/g11n/test"
)

func testMessage(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Message is not the same as expected.\n"+
			"\tActual: %v\n"+
			"\tExpected: %v\n", actual, expected)
	}
}

func TestSetLocaleFromRequest(t *testing.T) {
	type M struct {
		MyLittleSomething func() string `default:"cat"`
	}

	bgLocale := TempFile(`
	{
	  "M.MyLittleSomething": "котка"
	}
`)

	esLocale := TempFile(`
	{
	  "M.MyLittleSomething": "gato"
	}
`)

	factory := New()

	factory.SetLocales(map[language.Tag]string{
		language.Bulgarian: bgLocale,
		language.Spanish:   esLocale,
	}, "json")

	r, _ := http.NewRequest("GET", "https://golang.org", nil)
	r.Header.Add("Accept-Language", "bg")

	SetLocale(factory, r)

	m := factory.Init(&M{}).(*M)

	testMessage(t,
		string(m.MyLittleSomething()),
		`котка`)
}

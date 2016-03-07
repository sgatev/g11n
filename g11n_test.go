package g11n_test

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/text/language"

	. "github.com/s2gatev/g11n"
	. "github.com/s2gatev/g11n/test"
)

func testMessage(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Message is not the same as expected.\n"+
			"\tActual: %v\n"+
			"\tExpected: %v\n", actual, expected)
	}
}

func TestInitVar(t *testing.T) {
	type M struct {
		MyLittleSomething func() string `default:"Not as quick as the brown fox."`
	}

	var m M
	New().Init(&m)

	testMessage(t,
		m.MyLittleSomething(),
		"Not as quick as the brown fox.")
}

func TestStringReturnType(t *testing.T) {
	type M struct {
		MyLittleSomething string `default:"Not as quick as the brown fox."`
	}

	m := New().Init(&M{}).(*M)

	testMessage(t,
		m.MyLittleSomething,
		"Not as quick as the brown fox.")
}

func TestCustomReturnType(t *testing.T) {
	type M struct {
		MyLittleSomething SafeHTMLFormat `default:"<message>Oops!</message>"`
	}

	m := New().Init(&M{}).(*M)

	testMessage(t,
		string(m.MyLittleSomething),
		`\<message\>Oops!\<\/message\>`)
}

func TestSimpleMessage(t *testing.T) {
	type M struct {
		MyLittleSomething func() string `default:"Not as quick as the brown fox."`
	}

	m := New().Init(&M{}).(*M)

	testMessage(t,
		m.MyLittleSomething(),
		"Not as quick as the brown fox.")
}

func TestInitEmbeddedStruct(t *testing.T) {
	type N struct {
		MyLittleSomething func() string `default:"Not as quick as the brown fox."`
	}

	type M struct {
		*N
	}

	m := New().Init(&M{}).(*M)

	testMessage(t,
		m.MyLittleSomething(),
		"Not as quick as the brown fox.")
}

func TestMessageWithNumberArguments(t *testing.T) {
	type M struct {
		MyLittleSomething func(int, float64) string `default:"And yeah, it works: %v %v"`
	}

	m := New().Init(&M{}).(*M)

	testMessage(t,
		m.MyLittleSomething(42, 3.14),
		"And yeah, it works: 42 3.14")
}

func TestMessageWithCustomFormat(t *testing.T) {
	type M struct {
		MyLittleSomething func(CustomFormat) string `default:"Surprise! %v"`
	}

	m := New().Init(&M{}).(*M)

	testMessage(t,
		m.MyLittleSomething(CustomFormat{func() string {
			return "<ops>This works</ops>"
		}}),
		"Surprise! <ops>This works</ops>")
}

func TestPluralMessage(t *testing.T) {
	type M struct {
		MyLittleSomething func(PluralFormat) string `default:"Count: %v"`
	}

	m := New().Init(&M{}).(*M)

	testMessage(t,
		m.MyLittleSomething(0),
		"Count: some")
	testMessage(t,
		m.MyLittleSomething(1),
		"Count: crazy")
	testMessage(t,
		m.MyLittleSomething(21),
		"Count: stuff")
}

func TestMessageWithDifferentResult(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHTMLFormat `default:"<message>Oops!</message>"`
	}

	m := New().Init(&M{}).(*M)

	testMessage(t,
		string(m.MyLittleSomething()),
		`\<message\>Oops!\<\/message\>`)
}

func TestMessageWithMultipleResults(t *testing.T) {
	type M struct {
		MyLittleSomething func() (string, int) `default:"Oops!"`
	}

	defer MustPanic(t, "Wrong number of results in a g11n message. Expected 1, got 2.")

	New().Init(&M{})
}

func TestLocalizedMessageSingle(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHTMLFormat `default:"Cat"`
	}

	bgLocale := TempFile(`
	{
	  "M.MyLittleSomething": "Котка"
	}
`)

	factory := New()

	factory.SetLocale(language.Bulgarian, "json", bgLocale)
	factory.LoadLocale(language.Bulgarian)

	m := factory.Init(&M{}).(*M)

	testMessage(t,
		string(m.MyLittleSomething()),
		`Котка`)
}

func TestLocalizedMessageMultiple(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHTMLFormat `default:"Cat"`
	}

	bgLocale := TempFile(`
	{
	  "M.MyLittleSomething": "Котка"
	}
`)

	factory := New()

	factory.SetLocales(map[language.Tag]string{
		language.Bulgarian: bgLocale,
	}, "json")
	factory.LoadLocale(language.Bulgarian)

	m := factory.Init(&M{}).(*M)

	testMessage(t,
		string(m.MyLittleSomething()),
		`Котка`)
}

func TestSetLocaleAfterInit(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHTMLFormat `default:"Cat"`
	}

	bgLocale := TempFile(`
	{
	  "M.MyLittleSomething": "Котка"
	}
`)

	factory := New()

	m := factory.Init(&M{}).(*M)

	factory.SetLocale(language.Bulgarian, "json", bgLocale)
	factory.LoadLocale(language.Bulgarian)

	testMessage(t,
		string(m.MyLittleSomething()),
		`Котка`)
}

func TestLocalizedMessageUnknownFormat(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHTMLFormat `default:"Cat"`
	}

	bgLocale := TempFile(`
	M.MyLittleSomething: Котка
`)

	defer MustPanic(t, "Unknown locale format 'custom'.")

	factory := New()
	factory.SetLocale(language.Bulgarian, "custom", bgLocale)
	factory.LoadLocale(language.Bulgarian)
}

func TestLocales(t *testing.T) {
	factory := New()

	factory.SetLocales(map[language.Tag]string{
		language.Bulgarian: "",
		language.Spanish:   "",
		language.Italian:   "",
	}, "custom")

	expected := map[language.Tag]struct{}{
		language.Bulgarian: struct{}{},
		language.Spanish:   struct{}{},
		language.Italian:   struct{}{},
	}
	actual := map[language.Tag]struct{}{}
	for _, tag := range factory.Locales() {
		actual[tag] = struct{}{}
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Locales are not correct.\n"+
			"Expected: %v\n"+
			"Actual: %v\n", expected, actual)
	}
}

func TestUnknownLocale(t *testing.T) {
	defer MustPanic(t, "Unknown locale 'bg'.")

	factory := New()
	factory.LoadLocale(language.Bulgarian)
}

type CustomFormat struct {
	message func() string
}

func (cf CustomFormat) G11nParam() string {
	return cf.message()
}

type PluralFormat int

func (pf PluralFormat) G11nParam() string {
	switch pf {
	case 0:
		return "some"
	case 1:
		return "crazy"
	default:
		return "stuff"
	}
}

type SafeHTMLFormat string

func (shf SafeHTMLFormat) G11nResult(formattedMessage string) string {
	r := strings.NewReplacer("<", `\<`, ">", `\>`, "/", `\/`)
	return r.Replace(formattedMessage)
}

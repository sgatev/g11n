package g11n_test

import (
	"strings"
	"testing"

	. "github.com/s2gatev/g11n"
)

func TestSimpleMessage(t *testing.T) {
	type M struct {
		MyLittleSomething func() string `default:"Not as quick as the brown fox."`
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

	defer testPanic(t, "Wrong number of results in a g11n message. Expected 1, got 2.")

	New().Init(&M{})
}

func TestLocalizedMessage(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHTMLFormat `default:"Cat"`
	}

	bgLocale := tempFile(`
	{
	  "M.MyLittleSomething": "Котка"
	}
`)

	factory := New()

	factory.LoadLocale("json", "bg", bgLocale)
	factory.SetLocale("bg")

	m := factory.Init(&M{}).(*M)

	testMessage(t,
		string(m.MyLittleSomething()),
		`Котка`)
}

func TestLocalizedMessageUnknownFormat(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHTMLFormat `default:"Cat"`
	}

	bgLocale := tempFile(`
	M.MyLittleSomething: Котка
`)

	defer testPanic(t, "Unknown locale format 'custom'.")

	New().LoadLocale("custom", "bg", bgLocale)
}

func TestInitAsync(t *testing.T) {
	type M struct {
		MyLittleSomething func() string `default:"Not as quick as the brown fox."`
	}

	a, s := New().InitAsync(&M{})
	m := a.(*M)

	s.Await()

	testCompleted(t, s, true)
	testMessage(t,
		m.MyLittleSomething(),
		"Not as quick as the brown fox.")
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

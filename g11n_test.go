package g11n_test

import (
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/s2gatev/g11n"
)

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

type SafeHtmlFormat string

func (shf SafeHtmlFormat) G11nResult(formattedMessage string) string {
	r := strings.NewReplacer("<", `\<`, ">", `\>`, "/", `\/`)
	return r.Replace(formattedMessage)
}

func testStringsEqual(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Expected strings to be equal, buuut...\n"+
			"\tActual: %v\n"+
			"\tExpected: %v\n", actual, expected)
	}
}

func testPanic(t *testing.T, expectedMessage string) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic.")
	} else {
		actualMessage := r.(string)
		if expectedMessage != actualMessage {
			t.Errorf("The panic message was not correct.\n"+
				"\tExpected: %v\n"+
				"\tActual: %v\n", expectedMessage, actualMessage)
		}
	}
}

func tempFile(content string) string {
	file, err := ioutil.TempFile("", "g11n")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(file.Name(), []byte(content), 0644)
	return file.Name()
}

func TestSimpleMessage(t *testing.T) {
	type M struct {
		MyLittleSomething func() string `default:"Not as quick as the brown fox."`
	}

	m := New().Init(&M{}).(*M)

	testStringsEqual(t,
		m.MyLittleSomething(),
		"Not as quick as the brown fox.")
}

func TestMessageWithNumberArguments(t *testing.T) {
	type M struct {
		MyLittleSomething func(int, float64) string `default:"And yeah, it works: %v %v"`
	}

	m := New().Init(&M{}).(*M)

	testStringsEqual(t,
		m.MyLittleSomething(42, 3.14),
		"And yeah, it works: 42 3.14")
}

func TestMessageWithCustomFormat(t *testing.T) {
	type M struct {
		MyLittleSomething func(CustomFormat) string `default:"Surprise! %v"`
	}

	m := New().Init(&M{}).(*M)

	testStringsEqual(t,
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

	testStringsEqual(t,
		m.MyLittleSomething(0),
		"Count: some")
	testStringsEqual(t,
		m.MyLittleSomething(1),
		"Count: crazy")
	testStringsEqual(t,
		m.MyLittleSomething(21),
		"Count: stuff")
}

func TestMessageWithDifferentResult(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHtmlFormat `default:"<message>Oops!</message>"`
	}

	m := New().Init(&M{}).(*M)

	testStringsEqual(t,
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
		MyLittleSomething func() SafeHtmlFormat `default:"Cat"`
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

	testStringsEqual(t,
		string(m.MyLittleSomething()),
		`Котка`)
}

func TestLocalizedMessageUnknownFormat(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHtmlFormat `default:"Cat"`
	}

	bgLocale := tempFile(`
	M.MyLittleSomething: Котка
`)

	defer testPanic(t, "Unknown locale format 'custom'.")

	New().LoadLocale("custom", "bg", bgLocale)
}

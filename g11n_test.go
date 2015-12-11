package g11n_test

import (
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

func TestEmbedSimpleMessage(t *testing.T) {
	type M struct {
		MyLittleSomething func() string `default:"Not as quick as the brown fox."`
	}

	m := Init(&M{}).(*M)

	testStringsEqual(t,
		m.MyLittleSomething(),
		"Not as quick as the brown fox.")
}

func TestEmbedMessageWithNumberArguments(t *testing.T) {
	type M struct {
		MyLittleSomething func(int, float64) string `default:"And yeah, it works: %v %v"`
	}

	m := Init(&M{}).(*M)

	testStringsEqual(t,
		m.MyLittleSomething(42, 3.14),
		"And yeah, it works: 42 3.14")
}

func TestEmbedMessageWithCustomFormat(t *testing.T) {
	type M struct {
		MyLittleSomething func(CustomFormat) string `default:"Surprise! %v"`
	}

	m := Init(&M{}).(*M)

	testStringsEqual(t,
		m.MyLittleSomething(CustomFormat{func() string {
			return "<ops>This works</ops>"
		}}),
		"Surprise! <ops>This works</ops>")
}

func TestEmbedPluralMessage(t *testing.T) {
	type M struct {
		MyLittleSomething func(PluralFormat) string `default:"Count: %v"`
	}

	m := Init(&M{}).(*M)

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

func TestEmbedMessageWithDifferentResult(t *testing.T) {
	type M struct {
		MyLittleSomething func() SafeHtmlFormat `default:"<message>Oops!</message>"`
	}

	m := Init(&M{}).(*M)

	testStringsEqual(t,
		string(m.MyLittleSomething()),
		`\<message\>Oops!\<\/message\>`)
}

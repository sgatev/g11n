package g11n_test

import (
	"testing"

	. "github.com/s2gatev/g11n"
)

func testStringsEqual(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Expected strings to be equal, buuut...\n"+
			"\tActual: %v\n"+
			"\tExpected: %v\n", actual, expected)
	}
}

func TestEmbedSimpleMessage(t *testing.T) {
	type M struct {
		MyLittleSomething func() string `embed:"Not as quick as the brown fox."`
	}

	m := Init(&M{}).(*M)

	testStringsEqual(t,
		m.MyLittleSomething(),
		"Not as quick as the brown fox.")
}

func TestEmbedMessageWithNumberArguments(t *testing.T) {
	type M struct {
		MyLittleSomething func(int, float64) string `embed:"And yeah, it works: %v %v"`
	}

	m := Init(&M{}).(*M)

	testStringsEqual(t,
		m.MyLittleSomething(42, 3.14),
		"And yeah, it works: 42 3.14")
}

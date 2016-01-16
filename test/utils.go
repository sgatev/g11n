package test

import (
	"io/ioutil"
	"testing"
)

// MustPanic expects a panic in the current goroutine and verifies its
// message against an expected one.
func MustPanic(t *testing.T, expectedMessage string) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic.")
	} else {
		if actualMessage := r.(string); expectedMessage != actualMessage {
			t.Errorf("The panic message was not correct.\n"+
				"\tExpected: %v\n"+
				"\tActual: %v\n", expectedMessage, actualMessage)
		}
	}
}

// TempFile creates a temporary file with some content and returns the file name.
func TempFile(content string) string {
	file, err := ioutil.TempFile("", "g11n")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(file.Name(), []byte(content), 0644)

	return file.Name()
}

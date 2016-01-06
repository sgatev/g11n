package g11n_test

import (
	"io/ioutil"
	"testing"

	. "github.com/s2gatev/g11n"
)

// testCompleted verifies the completion status of a Synchronizer against
// an expected value.
func testCompleted(t *testing.T, synchronizer *Synchronizer, expected bool) {
	if actual := synchronizer.Completed(); actual != expected {
		t.Errorf("Asynchronous initialization status is not the same as expected.\n"+
			"\tActual: %v\n"+
			"\tExpected: %v\n", actual, expected)
	}
}

// testMessage verifies the string value of a message against an expected one.
func testMessage(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Message is not the same as expected.\n"+
			"\tActual: %v\n"+
			"\tExpected: %v\n", actual, expected)
	}
}

// testPanic expects a panic in the current goroutine and verifies its
// message against an expected one.
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

// tempFile creates a temporary file with some content and returns the file name.
func tempFile(content string) string {
	file, err := ioutil.TempFile("", "g11n")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(file.Name(), []byte(content), 0644)
	return file.Name()
}

package locale_test

import (
	"reflect"
	"testing"

	. "github.com/s2gatev/g11n/locale"
	. "github.com/s2gatev/g11n/test"
)

func testLoadJson(t *testing.T, filePath string, expected map[string]string) {
	loader, _ := GetLoader("json")

	if actual := loader.Load(filePath); !reflect.DeepEqual(actual, expected) {
		t.Errorf("")
	}
}

func TestLoadCorrectJson(t *testing.T) {
	filePath := TempFile(`
	{
	  "M.MyLittleSomething": "Котка"
	}
`)

	testLoadJson(t, filePath, map[string]string{
		"M.MyLittleSomething": "Котка",
	})
}

func TestLoadIncorrectJson(t *testing.T) {
	filePath := TempFile(`

	  "M.MyLittleSomething": "Котка"
	}
`)

	testLoadJson(t, filePath, map[string]string{})
}

func TestLoadJsonWithDuplicateKeys(t *testing.T) {
	filePath := TempFile(`
	{
	  "M.MyLittleSomething": "First",
		"M.MyLittleSomething": "Second"
	}
`)

	testLoadJson(t, filePath, map[string]string{
		"M.MyLittleSomething": "Second",
	})
}

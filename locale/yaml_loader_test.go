package locale_test

import (
	"reflect"
	"testing"

	. "github.com/s2gatev/g11n/locale"
	. "github.com/s2gatev/g11n/test"
)

func testLoadYaml(t *testing.T, filePath string, expected map[string]string) {
	loader, _ := GetLoader("yaml")

	if actual := loader.Load(filePath); !reflect.DeepEqual(actual, expected) {
		t.Errorf("")
	}
}

func TestLoadCorrectYaml(t *testing.T) {
	filePath := TempFile(`
M.MyLittleSomething: Котка
`)

	testLoadYaml(t, filePath, map[string]string{
		"M.MyLittleSomething": "Котка",
	})
}

func TestLoadIncorrectYaml(t *testing.T) {
	filePath := TempFile(`
M.MyLittleSomething - Котка
`)

	testLoadYaml(t, filePath, map[string]string{})
}

func TestLoadYamlWithDuplicateKeys(t *testing.T) {
	filePath := TempFile(`
M.MyLittleSomething: First
M.MyLittleSomething: Second
`)

	testLoadYaml(t, filePath, map[string]string{
		"M.MyLittleSomething": "Second",
	})
}

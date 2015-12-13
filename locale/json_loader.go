package locale

import (
	"encoding/json"
	"io/ioutil"
)

type jsonLoader struct{}

func (jl *jsonLoader) Load(fileName string) map[string]string {
	result := map[string]string{}
	if data, err := ioutil.ReadFile(fileName); err == nil {
		json.Unmarshal(data, &result)
	}
	return result
}

func init() {
	RegisterLoader("json", &jsonLoader{})
}

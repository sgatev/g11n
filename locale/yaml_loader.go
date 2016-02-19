package locale

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type yamlLoader struct{}

func (yl *yamlLoader) Load(fileName string) map[string]string {
	result := map[string]string{}
	if data, err := ioutil.ReadFile(fileName); err == nil {
		yaml.Unmarshal(data, &result)
	}
	return result
}

func init() {
	RegisterLoader("yaml", &yamlLoader{})
}

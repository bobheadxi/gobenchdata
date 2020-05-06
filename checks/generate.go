package checks

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// GenerateConfig outputs configuration in the provided directory
func GenerateConfig(dir string) error {
	b, _ := yaml.Marshal(&Config{})
	return ioutil.WriteFile(defaultConfigPath(dir), b, os.ModePerm)
}

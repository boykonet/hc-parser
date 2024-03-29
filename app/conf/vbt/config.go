package vbt

import (
	"bytes"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Domain			string `yaml: domain`
	Site			string `yaml: site`
	Cursor			string `yaml: cursor`
	Cookie struct {
		December string `yaml: december`
		January  string `yaml: january`
	} `yaml: cookie`
	PropertiesPath string `yaml: propertiespath`
}

func ParseConfiguration(pathToFile string) (*Configuration, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	conf := &Configuration{}

	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, file); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buff.Bytes(), conf); err != nil {
		return nil, err
	}
	return conf, nil
}

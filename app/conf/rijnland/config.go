package rijnland

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Domain			string `yaml: domain`
	Site			string `yaml: site`
	Cursor			string `yaml: cursor`
	Host			string `yaml: host`
	Connection		string `yaml: connection`
	UserAgent		string `yaml: useragent`
	Cookie			string `yaml: cookie`
	PropertiesPath	string `yaml: propertiespath`
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
		return nil, fmt.Errorf("%v - %v", err, string(buff.Bytes()))
	}
	return conf, nil
}
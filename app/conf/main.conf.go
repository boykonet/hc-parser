package conf

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Min int `yaml: min`
	Max int `yaml: max`
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

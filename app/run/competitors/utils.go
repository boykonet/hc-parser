package competitors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func ToFile(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error create file: %v", err)
	}
	defer file.Close()

	reader := bytes.NewReader(data)

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("error copy file: %v", err)
	}
	return nil
}

func fromFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []byte{}, nil	  
		} else {
			return nil, fmt.Errorf("check if file exists: %v", err)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error open file")
	}
	return data, nil
}

func SaveToFile(path string, flats map[string]struct{}) error {
	var properties []string

	for key := range flats {
		properties = append(properties, key)
	}

	var p Properties
	p.Properties = properties

	data, err := yaml.Marshal(&p)
	if err != nil {
		return err
	}	
	return ToFile(path, data)
}

func GetFromFile(path string, flats map[string]struct{}) error {
	data, err := fromFile(path)
	if err != nil {
		return err
	}

	var p Properties
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	for _, val := range p.Properties {
		flats[val] = struct{}{}
	}
	return nil
}
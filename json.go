package mappath

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
)

// FromJson is a factory method to create a MapPath from JSON byte data
func FromJson(in []byte) (*MapPath, error) {
	var data interface{}
	err := json.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}
	switch data.(type) {
	case map[string]interface{}:
		return NewMapPath(data.(map[string]interface{})), nil
	}

	return nil, fmt.Errorf("Cannot JSON which is marshalled to %+v. Must be marshallable to map[string]interface {}", reflect.TypeOf(data))
}

// FromJsonFile is a factory method to create a MapPath from a JSON file
func FromJsonFile(file string) (*MapPath, error) {
	in, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return FromJson(in)
}

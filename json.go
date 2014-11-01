package mappath

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
)

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

func FromJsonFile(file string) (*MapPath, error) {
	in, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fmt.Printf("JSON: %s\n\n", string(in))
	return FromJson(in)
}

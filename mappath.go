package mappath

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/*
 * ------
 * Types
 * ------
 */

// Branch is a shorthand for the map-string structures we're working with
type Branch map[string]interface{}

// MapPath is the primary object type this package is about
type MapPath struct {
	root Branch
}

/*
 * ------
 * Errors
 * ------
 */

// NotFoundError is returned if a given path cannot be found
type NotFoundError string

func (err NotFoundError) Error() string {
	return "The path \"" + string(err) + "\" does not exist"
}

// InvalidTypeError is returned if a type getter (eg GetInt) is used but the
//	found type cannot be converted
type InvalidTypeError struct {
	source   interface{}
	expected string
}

func (err *InvalidTypeError) Error() string {
	return fmt.Sprintf("Could not cast %v into %s", reflect.TypeOf(err.source), err.expected)
}

/*
 * ------
 * MapPath methods
 * ------
 */

// NewMapPath creates is the primary constructor
func NewMapPath(root map[string]interface{}) *MapPath {
	return &MapPath{root: root}
}

// Get returns object found with given path
func (this *MapPath) Get(path string, fallback ...interface{}) (interface{}, error) {
	val, found := this.getBranch(strings.Split(path, "/"), this.root)
	if found {
		return val, nil
	} else if len(fallback) > 0 {
		return fallback[0], nil
	}
	return nil, NotFoundError(path)
}

// Has check whether the given path exists
func (this *MapPath) Has(path string) bool {
	_, ok := this.getBranch(strings.Split(path, "/"), this.root)
	return ok
}

// GetInt returns int value of path. If value cannot be parsed or converted to
// int then an InvalidTypeError is returned
func (this *MapPath) GetInt(path string, fallback ...int) (int, error) {
	var val interface{}
	var err error
	if len(fallback) > 0 {
		val, err = this.Get(path, fallback[0])
	} else {
		val, err = this.Get(path)
	}
	if err != nil {
		return 0, err
	}

	switch reflect.TypeOf(val).Kind() {

	case reflect.String:
		r, err := strconv.Atoi(val.(string))
		if err != nil {
			r, ferr := strconv.ParseFloat(val.(string), 64)
			if ferr == nil {
				return int(r), nil
			}
			return 0, err
		}
		return r, nil

	case reflect.Int:
		return val.(int), nil

	case reflect.Float64:
		return int(val.(float64)), nil
	}

	return 0, &InvalidTypeError{val, "int"}
}

// GetFloat returns float64 value of path. If value cannot be parsed or converted to
// float64 then an InvalidTypeError is returned
func (this *MapPath) GetFloat(path string, fallback ...float64) (float64, error) {
	var val interface{}
	var err error
	if len(fallback) > 0 {
		val, err = this.Get(path, fallback[0])
	} else {
		val, err = this.Get(path)
	}
	if err != nil {
		return 0.0, err
	}
	switch reflect.TypeOf(val).Kind() {

	case reflect.String:
		r, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return 0.0, err
		}
		return r, nil

	case reflect.Float64:
		return val.(float64), nil

	case reflect.Int:
		return float64(val.(int)), nil
	}

	return 0.0, &InvalidTypeError{val, "float64"}
}

// GetString returns string value of path. If value cannot be parsed or converted to
// string then an InvalidTypeError is returned
func (this *MapPath) GetString(path string, fallback ...string) (string, error) {
	var val interface{}
	var err error
	if len(fallback) > 0 {
		val, err = this.Get(path, fallback[0])
	} else {
		val, err = this.Get(path)
	}
	if err != nil {
		return "", err
	}
	switch reflect.TypeOf(val).Kind() {

	case reflect.String:
		return val.(string), nil

	case reflect.Float64:
		return fmt.Sprintf("%.9f", val.(float64)), nil

	case reflect.Int:
		return fmt.Sprintf("%d", val.(int)), nil

	}

	return "", &InvalidTypeError{val, "float64"}
}

// GetMap returns map[string]interface{} value of path. If value cannot be parsed or converted to
// map[string]interface{} then an InvalidTypeError is returned
func (this *MapPath) GetMap(path string, fallback ...map[string]interface{}) (map[string]interface{}, error) {
	var val interface{}
	var err error
	if len(fallback) > 0 {
		val, err = this.Get(path, fallback[0])
	} else {
		val, err = this.Get(path)
	}
	if err != nil {
		return nil, err
	}

	switch val.(type) {
	case map[string]interface{}:
		return val.(map[string]interface{}), nil
	}

	return nil, &InvalidTypeError{val, "map"}
}

// get
func (this *MapPath) getBranch(pathParts []string, current map[string]interface{}) (interface{}, bool) {
	name := pathParts[0]
	val, ok := current[name]
	if !ok {
		return nil, false
	}

	return this.getNext(pathParts, val)
}

func (this *MapPath) getArray(pathParts []string, current reflect.Value) (interface{}, bool) {
	idx, err := strconv.Atoi(pathParts[0])
	if err != nil || idx < 0 || idx >= current.Len() {
		return nil, false
	}

	return this.getNext(pathParts, current.Index(idx).Interface())
}

func (this *MapPath) getNext(pathParts []string, val interface{}) (interface{}, bool) {
	if len(pathParts) > 1 {
		t := reflect.TypeOf(val)
		switch t.Kind() {
		case reflect.Map:
			return this.getBranch(pathParts[1:], val.(map[string]interface{}))
		case reflect.Slice:
			return this.getArray(pathParts[1:], reflect.ValueOf(val))
		default:
			return nil, false
		}
	} else {
		return val, true
	}
}

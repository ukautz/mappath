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
// found type cannot be converted
type InvalidTypeError struct {
	source   interface{}
	expected string
}

func (err *InvalidTypeError) Error() string {
	return fmt.Sprintf("Could not cast %v into %s", reflect.TypeOf(err.source), err.expected)
}

// UnsupportedTypeError is returned if an unsupported type is used
type UnsupportedTypeError string

func (err UnsupportedTypeError) Error() string {
	return fmt.Sprintf("Type %s is not supported", string(err))
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

// Root returns underly root map
func (this *MapPath) Root() map[string]interface{} {
	return this.root
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

// GetSub return sub MapPath object representing a branch of the data tree
// It does not support a fallback
func (this *MapPath) GetSub(path string, fallback ...*MapPath) (*MapPath, error) {
	branch, err := this.GetMap(path)
	if err != nil {
		if _, notFound := err.(NotFoundError); notFound && len(fallback) > 0 {
			return fallback[0], nil
		}
		return nil, err
	}

	return NewMapPath(branch), nil
}

// GetArray returns nested array of provided kind. Fallback values are not supported.
// You should use the specialized methods (GetInts, GetStrings..) unless you know what you are doing.
func (this *MapPath) GetArray(refType reflect.Type, path string) (interface{}, bool, error) {
	val, err := this.Get(path)
	if err != nil {
		return nil, false, err
	} else if reflect.Slice != reflect.TypeOf(val).Kind() {
		return nil, false, &InvalidTypeError{val, "array"}
	}

	refVal := reflect.ValueOf(val)
	if refVal.Len() == 0 {
		return nil, false, nil
	}

	var result interface{}
	switch refType.Kind() {
	case reflect.Int:
		result = make([]int, refVal.Len())
		break
	case reflect.Float64:
		result = make([]float64, refVal.Len())
		break
	case reflect.String:
		result = make([]string, refVal.Len())
		break
	case reflect.Map:
		result = make([]map[string]interface{}, refVal.Len())
		break
	default:
		return nil, false, UnsupportedTypeError(refType.Kind().String())
	}
	refResult := reflect.ValueOf(result)

	for i := 0; i < refVal.Len(); i++ {
		val := refVal.Index(i)
		if refType.Kind() == val.Kind() {
			refResult.Index(i).Set(refVal.Index(i))
		} else {

			// must convert or parse item
			switch refType.Kind() {

			// expecting []int
			case reflect.Int:
				switch val.Kind() {
				case reflect.Float64:
					refResult.Index(i).Set(refVal.Index(i).Convert(refType))
					break
				case reflect.String:
					v, eint := strconv.Atoi(refVal.Index(i).String())
					if eint != nil {
						f, _ := strconv.ParseFloat(refVal.Index(i).String(), 64)
						v = int(f)
					}
					refResult.Index(i).Set(reflect.ValueOf(v))
					break
				default:
					return nil, false, &InvalidTypeError{val, fmt.Sprintf("[%d] %s", i, refType.Kind())}
				}
				break

				// expecting []float64
			case reflect.Float64:
				switch val.Kind() {
				case reflect.Int:
					refResult.Index(i).Set(refVal.Index(i).Convert(refType))
					break
				case reflect.String:
					v, _ := strconv.ParseFloat(refVal.Index(i).String(), 64)
					refResult.Index(i).Set(reflect.ValueOf(v))
					break
				default:
					return nil, false, &InvalidTypeError{val, fmt.Sprintf("[%d] %s", i, refType.Kind())}
				}
				break

				// expecting []string
			case reflect.String:
				switch val.Kind() {
				case reflect.Int:
					refResult.Index(i).Set(reflect.ValueOf(fmt.Sprintf("%d", val.Int())))
					break
				case reflect.Float64:
					refResult.Index(i).Set(reflect.ValueOf(fmt.Sprintf("%.9f", val.Float())))
					break
				default:
					return nil, false, &InvalidTypeError{val, fmt.Sprintf("[%d] %s", i, refType.Kind())}
				}
				break
			default:
				return nil, false, &InvalidTypeError{val, fmt.Sprintf("[%d] %s", i, refType.Kind())}
			}
		}
	}

	return result, true, nil
}

// Get
func (this *MapPath) GetInts(path string, fallback ...[]int) ([]int, error) {
	res, found, err := this.GetArray(reflect.TypeOf(int(0)), path)
	if err != nil {
		if _, ok := err.(NotFoundError); len(fallback) > 0 && ok {
			return fallback[0], nil
		}
		return nil, err
	} else if !found {
		return []int{}, nil
	}
	return res.([]int), nil
}

func (this *MapPath) GetFloats(path string, fallback ...[]float64) ([]float64, error) {
	res, found, err := this.GetArray(reflect.TypeOf(float64(0.0)), path)
	if err != nil {
		if _, ok := err.(NotFoundError); len(fallback) > 0 && ok {
			return fallback[0], nil
		}
		return nil, err
	} else if !found {
		return []float64{}, nil
	}
	return res.([]float64), nil
}

func (this *MapPath) GetStrings(path string, fallback ...[]string) ([]string, error) {
	res, found, err := this.GetArray(reflect.TypeOf(string("")), path)
	if err != nil {
		if _, ok := err.(NotFoundError); len(fallback) > 0 && ok {
			return fallback[0], nil
		}
		return nil, err
	} else if !found {
		return []string{}, nil
	}
	return res.([]string), nil
}

// GetMaps returns list of maps, i.e. nested maps in an array
func (this *MapPath) GetMaps(path string, fallback ...[]map[string]interface{}) ([]map[string]interface{}, error) {
	res, found, err := this.GetArray(reflect.TypeOf(map[string]interface{}{}), path)
	if err != nil {
		if _, ok := err.(NotFoundError); len(fallback) > 0 && ok {
			return fallback[0], nil
		}
		return nil, err
	} else if !found {
		return []map[string]interface{}{}, nil
	}
	return res.([]map[string]interface{}), nil
}

// GetSubs return slist of MapPath sub objects
//   // Structure
func (this *MapPath) GetSubs(path string, fallback ...[]*MapPath) ([]*MapPath, error) {
	res, found, err := this.GetArray(reflect.TypeOf(map[string]interface{}{}), path)
	if err != nil {
		if _, ok := err.(NotFoundError); len(fallback) > 0 && ok {
			return fallback[0], nil
		}
		return nil, err
	} else if !found {
		return []*MapPath{}, nil
	}
	subs := make([]*MapPath, len(res.([]map[string]interface{})))
	for i, m := range res.([]map[string]interface{}) {
		subs[i] = &MapPath{m}
	}
	return subs, nil
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

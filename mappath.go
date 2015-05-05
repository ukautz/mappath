package mappath

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)


var kindsString = []reflect.Kind{
	reflect.String,
}
var kindsInt = []reflect.Kind{
	reflect.Int,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Int8,
	reflect.Uint,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Uint8,
}
var kindsFloat = []reflect.Kind{
	reflect.Float64,
	reflect.Float32,
}

func isOfKind(is reflect.Kind, anyOf []reflect.Kind) bool {
	for _, c := range anyOf {
		if is == c {
			return true
		}
	}
	return false
}

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

func (this *MapPath) GetAs(path string, typ reflect.Type, fallback ...interface{}) (interface{}, error) {
	val, err := this.Get(path, fallback...)
	if err != nil {
		return nil, err
	}
	kind := typ.Kind()
	valRef := reflect.ValueOf(val)
	valKind := valRef.Kind()

	switch {
		case isOfKind(kind, kindsString):
			switch {
				case isOfKind(valKind, kindsString):
					return val, nil
				case isOfKind(valKind, kindsInt):
					return fmt.Sprintf("%d", val), nil
				case isOfKind(valKind, kindsFloat):
					return fmt.Sprintf("%f", val), nil
				default:
					return fmt.Sprintf("%v", val), nil
			}
		case isOfKind(kind, kindsInt):
			switch {
				case isOfKind(valKind, kindsString):
					p, err := strconv.Atoi(val.(string))
					return p, err
				case isOfKind(valKind, kindsInt):
					return valRef.Convert(typ).Interface(), nil
				case isOfKind(valKind, kindsFloat):
					return valRef.Convert(typ).Interface(), nil
				default:
					return 0, &InvalidTypeError{val, "int"}
			}
		case isOfKind(kind, kindsFloat):
			switch {
				case isOfKind(valKind, kindsString):
					p, err := strconv.ParseFloat(val.(string), 64)
					return p, err
				case isOfKind(valKind, kindsInt):
					return valRef.Convert(typ).Interface(), nil
				case isOfKind(valKind, kindsFloat):
					return valRef.Convert(typ).Interface(), nil
				default:
					return 0.0, &InvalidTypeError{val, "float64"}
			}
		default:
			return nil, &InvalidTypeError{val, strings.ToLower(kind.String())}
	}
}

// Has check whether the given path exists
func (this *MapPath) Has(path string) bool {
	_, ok := this.getBranch(strings.Split(path, "/"), this.root)
	return ok
}

// GetInt returns int value of path. If value cannot be parsed or converted then an InvalidTypeError is returned
func (this *MapPath) Bool(path string, fallback ...bool) (bool, error) {
	var val interface{}
	var err error
	if len(fallback) > 0 {
		val, err = this.Get(path, fallback[0])
	} else {
		val, err = this.Get(path)
	}
	if err != nil {
		return false, err
	}
	switch reflect.TypeOf(val).Kind() {

		case reflect.Bool:
			return val.(bool), nil

		case reflect.Int:
			if val.(int) == 0 {
				return false, nil
			} else {
				return true, nil
			}

		case reflect.Float64:
			if val.(float64) == 0.0 {
				return false, nil
			} else {
				return true, nil
			}

		case reflect.String:
			switch val.(string) {
				case "true":
					return true, nil
				case "yes":
					return true, nil
				case "false":
					return false, nil
				case "no":
					return false, nil
				default:
					return false, fmt.Errorf("Cannot convert \"%s\" to bool (must be \"true\", \"yes\", \"false\" or \"no\")", val.(string))
			}
	}

	return false, &InvalidTypeError{val, "bool"}
}

// GetBoolV returns bool value of path. If value cannot be parsed or converted then fallback or false is returned. Handy in single value context.
func (this *MapPath) BoolV(path string, fallback ...bool) bool {
	if val, err := this.Bool(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return false
		}
	} else {
		return val
	}
}

// GetInt returns int value of path. If value cannot be parsed or converted then an InvalidTypeError is returned
func (this *MapPath) Int(path string, fallback ...int) (int, error) {
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
		case reflect.Bool:
			r := val.(bool)
			if r {
				return 1, nil
			} else {
				return 0, nil
			}

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

// GetIntV returns int value of path. If value cannot be parsed or converted then fallback or 0 is returned. Handy in single value context.
func (this *MapPath) IntV(path string, fallback ...int) int {
	if val, err := this.Int(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return 0
		}
	} else {
		return val
	}
}

// GetFloat returns float64 value of path. If value cannot be parsed or converted then an InvalidTypeError is returned
func (this *MapPath) Float(path string, fallback ...float64) (float64, error) {
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

		case reflect.Bool:
			r := val.(bool)
			if r {
				return 1.0, nil
			} else {
				return 0.0, nil
			}

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

// GetFloatV returns float64 value of path. If value cannot be parsed or converted then fallback or 0.0 is returned. Handy in single value context.
func (this *MapPath) FloatV(path string, fallback ...float64) float64 {
	if val, err := this.Float(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return 0.0
		}
	} else {
		return val
	}
}

// GetString returns string value of path. If value cannot be converted then an InvalidTypeError is returned
func (this *MapPath) String(path string, fallback ...string) (string, error) {
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

		case reflect.Bool:
			if val.(bool) {
				return "true", nil
			} else {
				return "false", nil
			}

		case reflect.String:
			return val.(string), nil

		case reflect.Float64:
			return fmt.Sprintf("%.9f", val.(float64)), nil

		case reflect.Int:
			return fmt.Sprintf("%d", val.(int)), nil

	}

	return "", &InvalidTypeError{val, "float64"}
}

// GetStringV returns string value of path. If value cannot be parsed or converted then fallback or "" is returned. Handy in single value context.
func (this *MapPath) StringV(path string, fallback ...string) string {
	if val, err := this.String(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return ""
		}
	} else {
		return val
	}
}

// GetMap returns the map value of path. If value is not a map then an InvalidTypeError is returned
func (this *MapPath) Map(path string, fallback ...map[string]interface{}) (map[string]interface{}, error) {
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
		case map[interface{}]interface{}:
			m := make(map[string]interface{})
			for k, v := range val.(map[interface{}]interface{}) {
				m[k.(string)] = v
			}
			return m, nil
	}

	return nil, &InvalidTypeError{val, "map"}
}

// GetMapV returns map[string]interface{} value of path. If value cannot be parsed or converted then fallback or nil is returned. Handy in single value context.
func (this *MapPath) MapV(path string, fallback ...map[string]interface{}) map[string]interface{} {
	if val, err := this.Map(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return nil
		}
	} else {
		return val
	}
}

// GetSub return a new MapPath object representing the sub structure, which needs to be a map. If the sub structure
// is of any other type then an InvalidTypeError is returned
func (this *MapPath) Child(path string, fallback ...*MapPath) (*MapPath, error) {
	branch, err := this.Map(path)
	if err != nil {
		if _, notFound := err.(NotFoundError); notFound && len(fallback) > 0 {
			return fallback[0], nil
		}
		return nil, err
	}

	return NewMapPath(branch), nil
}

// GetMapV returns *MapPath value of path. If value cannot be parsed or converted then fallback or nil is returned. Handy in single value context.
func (this *MapPath) ChildV(path string, fallback ...*MapPath) *MapPath {
	if val, err := this.Child(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return nil
		}
	} else {
		return val
	}
}

// GetArray returns nested array of provided type. Fallback values are not supported.
// If the path value is not an array then an InvalidTypeError is returned.
// You should use the specialized methods (GetInts, GetStrings..) unless you know what you are doing.
func (this *MapPath) Array(refType reflect.Type, path string) (interface{}, bool, error) {
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
			return nil, false, UnsupportedTypeError(refType.Kind().String()+ "@1")
	}
	refResult := reflect.ValueOf(result)

	for i := 0; i < refVal.Len(); i++ {
		itemRef := refVal.Index(i)
		if itemRef.Kind() == reflect.Interface {
			itemRef = reflect.ValueOf(itemRef.Interface())
		}
		if refType.Kind() == itemRef.Kind() && refType.ConvertibleTo(itemRef.Type()) {
			refResult.Index(i).Set(itemRef)
		} else {

			// must convert or parse item
			switch refType.Kind() {

				// expecting []int
				case reflect.Int:
					switch itemRef.Kind() {
						case reflect.Bool:
							if refVal.Index(i).Bool() {
								refResult.Index(i).Set(reflect.ValueOf(1))
							} else {
								refResult.Index(i).Set(reflect.ValueOf(0))
							}
							break
						case reflect.Float64:
							refResult.Index(i).Set(itemRef.Convert(refType))
							break
						case reflect.String:
							v, eint := strconv.Atoi(itemRef.String())
							if eint != nil {
								f, _ := strconv.ParseFloat(itemRef.String(), 64)
								v = int(f)
							}
							refResult.Index(i).Set(reflect.ValueOf(v))
							break
						default:
							return nil, false, &InvalidTypeError{itemRef.Interface(), fmt.Sprintf("[%d]array<%s>@2", i, refType.Kind())}
					}
					break

					// expecting []float64
				case reflect.Float64:
					switch itemRef.Kind() {
						case reflect.Bool:
							if refVal.Index(i).Bool() {
								refResult.Index(i).Set(reflect.ValueOf(1.0))
							} else {
								refResult.Index(i).Set(reflect.ValueOf(0.0))
							}
							break
						case reflect.Int:
							refResult.Index(i).Set(itemRef.Convert(refType))
							break
						case reflect.String:
							v, _ := strconv.ParseFloat(itemRef.String(), 64)
							refResult.Index(i).Set(reflect.ValueOf(v))
							break
						default:
							return nil, false, &InvalidTypeError{itemRef.Interface(), fmt.Sprintf("[%d]array<%s>@3", i, refType.Kind())}
						}
					break

					// expecting []string
				case reflect.String:
					switch itemRef.Kind() {
						case reflect.Bool:
							if itemRef.Bool() {
								refResult.Index(i).Set(reflect.ValueOf("true"))
							} else {
								refResult.Index(i).Set(reflect.ValueOf("false"))
							}
						break
						case reflect.Int:
							refResult.Index(i).Set(reflect.ValueOf(fmt.Sprintf("%d", itemRef.Int())))
							break
						case reflect.Float64:
							refResult.Index(i).Set(reflect.ValueOf(fmt.Sprintf("%.9f", itemRef.Float())))
							break
						case reflect.String:
							refResult.Index(i).Set(itemRef)
							break
						case reflect.Interface:
							s, ok := itemRef.Interface().(string)
							if !ok {
								return nil, false, &InvalidTypeError{itemRef.Interface(), fmt.Sprintf("[%d]array<%s>@4 - interface", i)}
							}
							refResult.Index(i).Set(reflect.ValueOf(s))
							break
						default:
							return nil, false, &InvalidTypeError{itemRef.Interface(), fmt.Sprintf("[%d]array<%s>@5 - %v", i, refType.Kind())}
					}
					break

					// expecting []map[string]interface{}
				case reflect.Map:
					var mapVal map[string]interface{}
					var ok bool
					if mapVal, ok = refVal.Index(i).Interface().(map[string]interface{}); !ok {
						if mapValRaw, ok := refVal.Index(i).Interface().(map[interface{}]interface{}); !ok {
							return nil, false, &InvalidTypeError{itemRef.Interface(), fmt.Sprintf("[%d]array<%s>@6", i, refType.Kind())}
						} else {
							mapVal = make(map[string]interface{})
							for k, v := range mapValRaw {
								mapVal[fmt.Sprintf("%s", k)] = v
							}
						}
					}
					if mapVal != nil {
						refResult.Index(i).Set(reflect.ValueOf(mapVal))
					}
					break

					// oops
				default:
					return nil, false, &InvalidTypeError{itemRef.Interface(), fmt.Sprintf("[%d]array<%s>@7", i, refType.Kind())}
			}
		}
	}

	return result, true, nil
}

// GetInts returns an array of int values. Tries to convert (eg float) or parse (string) values. If the
// path value cannot be parsed or converted than an InvalidTypeError is returned.
func (this *MapPath) Ints(path string, fallback ...[]int) ([]int, error) {
	res, found, err := this.Array(reflect.TypeOf(int(0)), path)
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

// GetIntsV returns []int value of path. If value cannot be parsed or converted then fallback or nil is returned. Handy in single value context.
func (this *MapPath) IntsV(path string, fallback ...[]int) []int {
	if val, err := this.Ints(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return nil
		}
	} else {
		return val
	}
}

// GetFloats returns an array of float64 values. Tries to convert (eg int) or parse (string) values. If the
// path value cannot be parsed or converted than an InvalidTypeError is returned.
func (this *MapPath) Floats(path string, fallback ...[]float64) ([]float64, error) {
	res, found, err := this.Array(reflect.TypeOf(float64(0.0)), path)
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

// GetFloatsV returns []float64 value of path. If value cannot be parsed or converted then fallback or nil is returned. Handy in single value context.
func (this *MapPath) FloatsV(path string, fallback ...[]float64) []float64 {
	if val, err := this.Floats(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return nil
		}
	} else {
		return val
	}
}

// GetStrings returns an array of string values. If the path value is incomaptible (eg map array) then an InvalidTypeError
// is returned
func (this *MapPath) Strings(path string, fallback ...[]string) ([]string, error) {
	res, found, err := this.Array(reflect.TypeOf(string("")), path)
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

// GetStringsV returns []string value of path. If value cannot be parsed or converted then fallback or nil is returned. Handy in single value context.
func (this *MapPath) StringsV(path string, fallback ...[]string) []string {
	if val, err := this.Strings(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return nil
		}
	} else {
		return val
	}
}

// GetMaps returns a nested array of maps. If the path value is not an array of maps then an InvalidTypeError is returned.
func (this *MapPath) Maps(path string, fallback ...[]map[string]interface{}) ([]map[string]interface{}, error) {
	res, found, err := this.Array(reflect.TypeOf(map[string]interface{}{}), path)
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

// GetMapsV returns []map[string]interface{} value of path. If value cannot be parsed or converted then fallback or nil is returned. Handy in single value context.
func (this *MapPath) MapsV(path string, fallback ...[]map[string]interface{}) []map[string]interface{} {
	if val, err := this.Maps(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return nil
		}
	} else {
		return val
	}
}

// GetSubs returns a nested array of sub structures. If the path value is not an array of maps then an InvalidTypeError is returned.
func (this *MapPath) Childs(path string, fallback ...[]*MapPath) ([]*MapPath, error) {
	res, found, err := this.Array(reflect.TypeOf(map[string]interface{}{}), path)
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

// GetSubsV returns []*MapPath value of path. If value cannot be parsed or converted then fallback or nil is returned. Handy in single value context.
func (this *MapPath) ChildsV(path string, fallback ...[]*MapPath) []*MapPath {
	if val, err := this.Childs(path, fallback...); err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		} else {
			return nil
		}
	} else {
		return val
	}
}

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
			m, ok := val.(map[string]interface{})
			if !ok {
				m = make(map[string]interface{})
				for k, v := range val.(map[interface{}]interface{}) {
					m[fmt.Sprintf("%s", k)] = v
				}
			}
			return this.getBranch(pathParts[1:], m)
		case reflect.Slice:
			return this.getArray(pathParts[1:], reflect.ValueOf(val))
		default:
			return nil, false
		}
	} else {
		return val, true
	}
}

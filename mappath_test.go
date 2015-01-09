package mappath

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var defaultTest = map[string]interface{}{
	"hello": "world",
	"bool": map[string]interface{}{
		"yes": true,
		"no": false,
		"stringyes1": "true",
		"stringyes2": "yes",
		"stringyes3": "notworking",
		"stringno1": "false",
		"stringno2": "no",
		"stringno3": "notworking",
	},
	"foo": map[string]interface{}{
		"bar": "baz",
		"baz": map[string]interface{}{
			"bam": 42,
		},
	},
	"array": map[string]interface{}{
		"empty":        []interface{}{},
		"realints":     []int{1, 2, 3, 4},
		"realfloats":   []float64{1.01, 2.02, 3.03, 4.04},
		"realbools":    []bool{true, true, false, false},
		"stringints":   []string{"1", "2", "3", "4"},
		"stringfloats": []string{"1.01", "2.02", "3.03", "4.04"},
		"stringbools":  []string{"true", "yes", "false", "no"},
		"strings":      []string{"foo", "bar", "baz"}},
	"3d-array": [][][]int{
		[][]int{
			[]int{1, 2, 3},
			[]int{4, 5, 6},
		},
		[][]int{
			[]int{11, 12, 13},
			[]int{14, 15, 16},
		},
	},
	"mixed": map[string]interface{}{
		"array1": []int{1, 2, 3, 4},
		"array2": []map[string]interface{}{
			map[string]interface{}{
				"foo": []int{1, 2, 3, 4},
				"bar": []string{"one", "two"},
			},
			map[string]interface{}{
				"foo": []int{11, 12, 13, 14},
				"bar": []string{"five", "six"},
			},
		},
	},
	"scalar": map[string]interface{}{
		"stringint":   "123",
		"stringfloat": "123.456",
		"realint":     123,
		"realfloat":   123.456,
	},
}

/*
 * -------
 * Root
 * -------
 */

var rootAccessTests = []map[string]interface{}{
	map[string]interface{}{
		"foo": "bar",
	},
	map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
	},
	map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"baz": []int{1, 2, 3},
			},
		},
	},
}

func TestRootAccess(t *testing.T) {
	for i, test := range rootAccessTests {
		mp := NewMapPath(test)
		assert.Equal(t, test, mp.Root(), fmt.Sprintf("Root %d kept", i))
	}
}

/*
 * -------
 * Get
 * -------
 */

var getExistingPathTests = []struct {
	path   string
	expect interface{}
	from   map[string]interface{}
}{

	// map access
	{
		path:   "hello",
		expect: "world",
		from:   defaultTest,
	},
	{
		path:   "foo/bar",
		expect: "baz",
		from:   defaultTest,
	},
	{
		path:   "foo/baz/bam",
		expect: 42,
		from:   defaultTest,
	},
	// array access
	{
		path:   "array/realints/0",
		expect: 1,
		from:   defaultTest,
	},
	{
		path:   "array/realints/3",
		expect: 4,
		from:   defaultTest,
	},
	{
		path:   "3d-array/0/0/0",
		expect: 1,
		from:   defaultTest,
	},
	{
		path:   "3d-array/1/0/0",
		expect: 11,
		from:   defaultTest,
	},
	{
		path:   "3d-array/1/1/2",
		expect: 16,
		from:   defaultTest,
	},
	// mixed access
	{
		path:   "mixed/array1/0",
		expect: 1,
		from:   defaultTest,
	},
	{
		path:   "mixed/array2/0/foo/0",
		expect: 1,
		from:   defaultTest,
	},
	{
		path:   "mixed/array2/0/bar/1",
		expect: "two",
		from:   defaultTest,
	},
	{
		path:   "mixed/array2/1/bar/1",
		expect: "six",
		from:   defaultTest,
	},
	// structure access
	{
		path: "foo/baz",
		expect: map[string]interface{}{
			"bam": 42,
		},
		from: defaultTest,
	},
	{
		path:   "array/realints",
		expect: []int{1, 2, 3, 4},
		from:   defaultTest,
	},
}

func TestGetExistingPath(t *testing.T) {
	for _, test := range getExistingPathTests {
		m := NewMapPath(test.from)
		r, e := m.Get(test.path)
		assert.Nil(t, e, "No error")
		assert.Equal(t, test.expect, r, fmt.Sprintf("Expected for %s: %+v", test.path, test.expect))
	}
}

func TestGetErrorOnMissingPath(t *testing.T) {
	for _, path := range []string{"foo", "foo/bar", "foo/bar/baz"} {
		m := NewMapPath(map[string]interface{}{})
		r, e := m.Get(path)
		assert.Nil(t, r, "Response is nil")
		assert.NotNil(t, e, "Error responded")
		assert.IsType(t, reflect.TypeOf(NotFoundError("")), reflect.TypeOf(e), "Correct error responded")
	}
}

func TestGetErrorOnWrongPath(t *testing.T) {
	for _, path := range []string{"bar", "foo/foo", "foo/bar/foo", "array/5", "3d-array/0/0/4", "3d-array/4/0/0"} {
		m := NewMapPath(defaultTest)
		r, e := m.Get(path)
		assert.Nil(t, r, "Response is nil")
		assert.NotNil(t, e, "Error responded")
		assert.IsType(t, reflect.TypeOf(NotFoundError("")), reflect.TypeOf(e), "Correct error responded")
	}
}

/*
 * -------
 * Has
 * -------
 */

func TestHasExistingPath(t *testing.T) {
	for _, test := range getExistingPathTests {
		m := NewMapPath(test.from)
		r := m.Has(test.path)
		assert.True(t, r, "Path found")
	}
}

func TestHasErrorOnMissingMapPath(t *testing.T) {
	for _, path := range []string{"foo", "foo/bar", "foo/bar/baz"} {
		m := NewMapPath(map[string]interface{}{})
		r := m.Has(path)
		assert.False(t, r, "Path not found")
	}
}

func TestHasErrorOnWrongMapPath(t *testing.T) {
	for _, path := range []string{"bar", "foo/foo", "foo/bar/foo", "array/5", "3d-array/0/0/4", "3d-array/4/0/0"} {
		m := NewMapPath(defaultTest)
		r := m.Has(path)
		assert.False(t, r, "Path not found")
	}
}

/*
 * -------
 * Get with fallback
 * -------
 */

var useGivenFallbackOnMissingPathTests = []struct {
	path     string
	fallback interface{}
}{
	// string fallback
	{
		path:     "foo",
		fallback: "bar",
	},
	// int fallback
	{
		path:     "foo",
		fallback: 42,
	},
	// deep path fallback
	{
		path:     "foo/bar/baz",
		fallback: "FOO!",
	},
}

func TestUseGivenFallbackOnMissingPath(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	for _, test := range useGivenFallbackOnMissingPathTests {
		r, e := m.Get(test.path, test.fallback)
		assert.Equal(t, r, test.fallback, "Fallback is returned")
		assert.Nil(t, e, "Error becomes nil")
	}
}

/*
 * -------
 * Get: Int
 * -------
 */

var getIntValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from actual int
	{
		path:     "scalar/realint",
		err:      false,
		expected: 123,
	},
	// from actual bool
	{
		path:     "bool/yes",
		err:      false,
		expected: 1,
	},
	{
		path:     "bool/no",
		err:      false,
		expected: 0,
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      false,
		expected: 123,
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      false,
		expected: 123,
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      false,
		expected: 123,
	},
	// from not parsable string
	{
		path:     "foo/bar",
		err:      true,
		expected: 0,
	},
	// from not parsable struct
	{
		path:     "foo/baz",
		err:      true,
		expected: 0,
	},
	// from not parsable array
	{
		path:     "array/realints",
		err:      true,
		expected: 0,
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: 0,
	},
}

func TestGetIntValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getIntValueTests {
		r, e := m.GetInt(test.path)
		if test.err {
			assert.NotNil(t, e, "Error returned OK")
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded")
		} else {
			assert.Nil(t, e, "NO error returned")
		}
		assert.Equal(t, test.expected, r, "Expected value returned")
	}
}

func TestGetIntValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := 345
	r, e := m.GetInt("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path")
	assert.Equal(t, r, f, "Fallback is returned")
}

func TestGetIntSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getIntValueTests {
		r := m.GetIntV(test.path)
		if test.err {
			assert.Equal(t, 0, r, "Nil value returned")
		}
		assert.Equal(t, test.expected, r, "Expected value returned")
	}
}

/*
 * -------
 * Get: Float
 * -------
 */

var getFloatValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from actual int
	{
		path:     "scalar/realint",
		err:      false,
		expected: 123.0,
	},
	// from actual bool
	{
		path:     "bool/yes",
		err:      false,
		expected: 1.0,
	},
	{
		path:     "bool/no",
		err:      false,
		expected: 0.0,
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      false,
		expected: 123.456,
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      false,
		expected: 123.0,
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      false,
		expected: 123.456,
	},
	// from not parsable string
	{
		path:     "foo/bar",
		err:      true,
		expected: 0.0,
	},
	// from not parsable struct
	{
		path:     "foo/baz",
		err:      true,
		expected: 0.0,
	},
	// from not parsable array
	{
		path:     "array/realints",
		err:      true,
		expected: 0.0,
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: 0.0,
	},
}

func TestGetFloatValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getFloatValueTests {
		r, e := m.GetFloat(test.path)
		if test.err {
			assert.NotNil(t, e, "Error returned OK")
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded")
		} else {
			assert.Nil(t, e, "NO error returned")
		}
		assert.Equal(t, test.expected, r, "Expected value returned")
	}
}

func TestGetFloatValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := 345.678
	r, e := m.GetFloat("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path")
	assert.Equal(t, r, f, "Fallback is returned")
}

func TestGetFloatSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getFloatValueTests {
		r := m.GetFloatV(test.path)
		if test.err {
			assert.Equal(t, 0.0, r, "Nil result returned")
		}
		assert.Equal(t, test.expected, r, "Expected value returned")
	}
}

/*
 * -------
 * Get: String
 * -------
 */

var getStringValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from actual int
	{
		path:     "scalar/realint",
		err:      false,
		expected: "123",
	},
	// from actual bool
	{
		path:     "bool/yes",
		err:      false,
		expected: "true",
	},
	{
		path:     "bool/no",
		err:      false,
		expected: "false",
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      false,
		expected: "123.456000000",
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      false,
		expected: "123",
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      false,
		expected: "123.456",
	},
	// from not regular string
	{
		path:     "foo/bar",
		err:      false,
		expected: "baz",
	},
	// from not parsable struct
	{
		path:     "foo/baz",
		err:      true,
		expected: "",
	},
	// from not parsable array
	{
		path:     "array/realints",
		err:      true,
		expected: "",
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: "",
	},
}

func TestGetStringValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getStringValueTests {
		r, e := m.GetString(test.path)
		if test.err {
			assert.NotNil(t, e, "Error returned OK")
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded")
		} else {
			assert.Nil(t, e, "NO error returned")
		}
		assert.Equal(t, test.expected, r, fmt.Sprintf("Expected value returned on %s", test.path))
	}
}

func TestGetStringValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := "foo"
	r, e := m.GetString("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path")
	assert.Equal(t, r, f, "Fallback is returned")
}

func TestGetStringSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getStringValueTests {
		r := m.GetStringV(test.path)
		if test.err {
			assert.Equal(t, "", r, "Nil result returned")
		}
		assert.Equal(t, test.expected, r, "Expected value returned")
	}
}

/*
 * -------
 * Get: Map
 * -------
 */

var getMapValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from actual int
	{
		path:     "scalar/realint",
		err:      true,
		expected: nil,
	},
	// from actual bool
	{
		path:     "bool/yes",
		err:      true,
		expected: nil,
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      true,
		expected: nil,
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      true,
		expected: nil,
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      true,
		expected: nil,
	},
	// from not regular string
	{
		path:     "foo/bar",
		err:      true,
		expected: nil,
	},
	// from parsable struct
	{
		path: "foo/baz",
		err:  false,
		expected: map[string]interface{}{
			"bam": 42,
		},
	},
	// from not parsable array
	{
		path:     "array/realints",
		err:      true,
		expected: nil,
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: nil,
	},
}

func TestGetMapValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getMapValueTests {
		r, e := m.GetMap(test.path)
		if test.err {
			assert.NotNil(t, e, fmt.Sprintf("Error has been returned on %s", test.path))
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded "+test.path)
		} else {
			assert.Nil(t, e, fmt.Sprintf("NO error returned on %s (%+v)", test.path, r))
		}
		if test.expected == nil {
			assert.Nil(t, r, fmt.Sprintf("Expected nil returned on %s", test.path))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("Expected value returned on %s", test.path))
		}
	}
}

func TestGetMapValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := map[string]interface{}{
		"foo": "bar",
	}
	r, e := m.GetMap("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path")
	assert.Equal(t, r, f, "Fallback is returned")
}

func TestGetMapSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for i, test := range getMapValueTests {
		r := m.GetMapV(test.path)
		if test.err {
			assert.Nil(t, r, fmt.Sprintf("[%d] Nil result returned", i))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("[%d] Expected value returned (ACTUAL: %+v)", i, r))
		}
	}
}

/*
 * -------
 * Get: Sub
 * -------
 */

var getSubValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from actual int
	{
		path:     "scalar/realint",
		err:      true,
		expected: nil,
	},
	// from actual bool
	{
		path:     "bool/yes",
		err:      true,
		expected: nil,
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      true,
		expected: nil,
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      true,
		expected: nil,
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      true,
		expected: nil,
	},
	// from not regular string
	{
		path:     "foo/bar",
		err:      true,
		expected: nil,
	},
	// from parsable struct
	{
		path: "foo/baz",
		err:  false,
		expected: NewMapPath(map[string]interface{}{
			"bam": 42,
		}),
	},
	// from not parsable array
	{
		path:     "array/realints",
		err:      true,
		expected: nil,
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: nil,
	},
}

func TestGetSubValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getSubValueTests {
		r, e := m.GetSub(test.path)
		if test.err {
			assert.NotNil(t, e, fmt.Sprintf("Error has been returned on %s", test.path))
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded "+test.path)
		} else {
			assert.Nil(t, e, fmt.Sprintf("NO error returned on %s (%+v)", test.path, r))
		}
		if test.expected == nil {
			assert.Nil(t, r, fmt.Sprintf("Expected nil returned on %s", test.path))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("Expected value returned on %s", test.path))
		}
	}
}

func TestGetSubValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := NewMapPath(map[string]interface{}{
		"foo": "bar",
	})
	r, e := m.GetSub("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path")
	assert.Equal(t, r, f, "Fallback is returned")
}

func TestGetSubSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for i, test := range getSubValueTests {
		r := m.GetSubV(test.path)
		if test.err {
			assert.Nil(t, r, fmt.Sprintf("[%d] Nil result returned", i))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("[%d] Expected value returned (ACTUAL: %+v)", i, r))
		}
	}
}

/*
 * -------
 * Get: Ints (list)
 * -------
 */

var getIntsValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from single int
	{
		path:     "scalar/realint",
		err:      true,
		expected: nil,
	},
	// from single bool
	{
		path:     "bool/yes",
		err:      true,
		expected: nil,
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      true,
		expected: nil,
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      true,
		expected: nil,
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      true,
		expected: nil,
	},
	// from not regular string
	{
		path:     "foo/bar",
		err:      true,
		expected: nil,
	},
	// from parsable struct
	{
		path:     "foo/baz",
		err:      true,
		expected: nil,
	},
	// from array of ints
	{
		path:     "array/realints",
		err:      false,
		expected: []int{1, 2, 3, 4},
	},
	// from array of bools
	{
		path:     "array/realbools",
		err:      false,
		expected: []int{1, 1, 0, 0},
	},
	// from array of floats
	{
		path:     "array/realfloats",
		err:      false,
		expected: []int{1, 2, 3, 4},
	},
	// from array of string ints
	{
		path:     "array/stringints",
		err:      false,
		expected: []int{1, 2, 3, 4},
	},
	// from array of string floats
	{
		path:     "array/stringfloats",
		err:      false,
		expected: []int{1, 2, 3, 4},
	},
	// from empty array
	{
		path:     "array/empty",
		err:      false,
		expected: []int{},
	},
	// from un-convertable array
	{
		path:     "mixed/array2",
		err:      true,
		expected: nil,
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: nil,
	},
}

func TestGetIntsValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getIntsValueTests {
		r, e := m.GetInts(test.path)
		if test.err {
			assert.NotNil(t, e, fmt.Sprintf("Error has been returned on %s", test.path))
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded "+test.path)
		} else {
			assert.Nil(t, e, fmt.Sprintf("NO error returned on %s (%+v)", test.path, r))
		}
		if test.expected == nil {
			assert.Nil(t, r, fmt.Sprintf("Expected nil returned on %s", test.path))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("Expected value returned on %s", test.path))
		}
	}
}

func TestGetIntsValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := []int{2, 3, 4}
	r, e := m.GetInts("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path (ints)")
	assert.Equal(t, r, f, "Fallback is returned (ints)")
}

func TestGetIntsSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for i, test := range getIntsValueTests {
		r := m.GetIntsV(test.path)
		if test.err {
			assert.Nil(t, r, fmt.Sprintf("[%d] Nil result returned", i))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("[%d] Expected value returned (ACTUAL: %+v)", i, r))
		}
	}
}

/*
 * -------
 * Get: Floats (list)
 * -------
 */

var getFloatsValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from single int
	{
		path:     "scalar/realint",
		err:      true,
		expected: nil,
	},
	// from single bool
	{
		path:     "bool/yes",
		err:      true,
		expected: nil,
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      true,
		expected: nil,
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      true,
		expected: nil,
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      true,
		expected: nil,
	},
	// from not regular string
	{
		path:     "foo/bar",
		err:      true,
		expected: nil,
	},
	// from parsable struct
	{
		path:     "foo/baz",
		err:      true,
		expected: nil,
	},
	// from array of ints
	{
		path:     "array/realints",
		err:      false,
		expected: []float64{1.0, 2.0, 3.0, 4.0},
	},
	// from array of bools
	{
		path:     "array/realbools",
		err:      false,
		expected: []float64{1.0, 1.0, 0.0, 0.0},
	},
	// from array of floats
	{
		path:     "array/realfloats",
		err:      false,
		expected: []float64{1.01, 2.02, 3.03, 4.04},
	},
	// from array of string ints
	{
		path:     "array/stringints",
		err:      false,
		expected: []float64{1.0, 2.0, 3.0, 4.0},
	},
	// from array of string floats
	{
		path:     "array/stringfloats",
		err:      false,
		expected: []float64{1.01, 2.02, 3.03, 4.04},
	},
	// from empty array
	{
		path:     "array/empty",
		err:      false,
		expected: []float64{},
	},
	// from un-convertable array
	{
		path:     "mixed/array2",
		err:      true,
		expected: nil,
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: nil,
	},
}

func TestGetFloatsValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getFloatsValueTests {
		r, e := m.GetFloats(test.path)
		if test.err {
			assert.NotNil(t, e, fmt.Sprintf("Error has been returned on %s", test.path))
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded "+test.path)
		} else {
			assert.Nil(t, e, fmt.Sprintf("NO error returned on %s (%+v)", test.path, r))
		}
		if test.expected == nil {
			assert.Nil(t, r, fmt.Sprintf("Expected nil returned on %s", test.path))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("Expected value returned on %s", test.path))
		}
	}
}

func TestGetFloatsValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := []float64{2.02, 3.03, 4.04}
	r, e := m.GetFloats("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path (floats)")
	assert.Equal(t, r, f, "Fallback is returned (floats)")
}

func TestGetFloatsSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for i, test := range getFloatsValueTests {
		r := m.GetFloatsV(test.path)
		if test.err {
			assert.Nil(t, r, fmt.Sprintf("[%d] Nil result returned", i))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("[%d] Expected value returned (ACTUAL: %+v)", i, r))
		}
	}
}

/*
 * -------
 * Get: Strings (list)
 * -------
 */

var getStringsValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from single int
	{
		path:     "scalar/realint",
		err:      true,
		expected: nil,
	},
	// from single bool
	{
		path:     "bool/yes",
		err:      true,
		expected: nil,
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      true,
		expected: nil,
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      true,
		expected: nil,
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      true,
		expected: nil,
	},
	// from not regular string
	{
		path:     "foo/bar",
		err:      true,
		expected: nil,
	},
	// from parsable struct
	{
		path:     "foo/baz",
		err:      true,
		expected: nil,
	},
	// from array of ints
	{
		path:     "array/realints",
		err:      false,
		expected: []string{"1", "2", "3", "4"},
	},
	// from array of bools
	{
		path:     "array/realbools",
		err:      false,
		expected: []string{"true", "true", "false", "false"},
	},
	// from array of floats
	{
		path:     "array/realfloats",
		err:      false,
		expected: []string{"1.010000000", "2.020000000", "3.030000000", "4.040000000"},
	},
	// from array of ints
	{
		path:     "array/stringints",
		err:      false,
		expected: []string{"1", "2", "3", "4"},
	},
	// from array of floats
	{
		path:     "array/stringfloats",
		err:      false,
		expected: []string{"1.01", "2.02", "3.03", "4.04"},
	},
	// from array of floats
	{
		path:     "array/strings",
		err:      false,
		expected: []string{"foo", "bar", "baz"},
	},
	// from empty array
	{
		path:     "array/empty",
		err:      false,
		expected: []string{},
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: nil,
	},
	// from invalid path
	{
		path:     "mixed/array2",
		err:      true,
		expected: nil,
	},
}

func TestGetStringsValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getStringsValueTests {
		r, e := m.GetStrings(test.path)
		if test.err {
			assert.NotNil(t, e, fmt.Sprintf("Error has been returned on %s", test.path))
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded "+test.path)
		} else {
			assert.Nil(t, e, fmt.Sprintf("NO error returned on %s (%+v)", test.path, r))
		}
		if test.expected == nil {
			assert.Nil(t, r, fmt.Sprintf("Expected nil returned on %s", test.path))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("Expected value returned on %s", test.path))
		}
	}
}

func TestGetStringsValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := []string{"a", "b"}
	r, e := m.GetStrings("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path (ints)")
	assert.Equal(t, r, f, "Fallback is returned (ints)")
}

func TestGetStringsSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for i, test := range getStringsValueTests {
		r := m.GetStringsV(test.path)
		if test.err {
			assert.Nil(t, r, fmt.Sprintf("[%d] Nil result returned", i))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("[%d] Expected value returned (ACTUAL: %+v)", i, r))
		}
	}
}

/*
 * -------
 * Get: Maps (list)
 * -------
 */

var getMapsValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from single int
	{
		path:     "scalar/realint",
		err:      true,
		expected: nil,
	},
	// from single bool
	{
		path:     "bool/yes",
		err:      true,
		expected: nil,
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      true,
		expected: nil,
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      true,
		expected: nil,
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      true,
		expected: nil,
	},
	// from not regular string
	{
		path:     "foo/bar",
		err:      true,
		expected: nil,
	},
	// from parsable struct
	{
		path:     "foo/baz",
		err:      true,
		expected: nil,
	},
	// from array of ints
	{
		path:     "array/realints",
		err:      true,
		expected: nil,
	},
	// from array of bools
	{
		path:     "array/realbools",
		err:      true,
		expected: nil,
	},
	// from array of floats
	{
		path:     "array/realfloats",
		err:      true,
		expected: nil,
	},
	// from array of ints
	{
		path:     "array/stringints",
		err:      true,
		expected: nil,
	},
	// from array of floats
	{
		path:     "array/stringfloats",
		err:      true,
		expected: nil,
	},
	// from array of floats
	{
		path:     "array/strings",
		err:      true,
		expected: nil,
	},
	// from empty array
	{
		path:     "array/empty",
		err:      false,
		expected: []map[string]interface{}{},
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: nil,
	},
	// from invalid path
	{
		path: "mixed/array2",
		err:  false,
		expected: []map[string]interface{}{
			map[string]interface{}{
				"foo": []int{1, 2, 3, 4},
				"bar": []string{"one", "two"},
			},
			map[string]interface{}{
				"foo": []int{11, 12, 13, 14},
				"bar": []string{"five", "six"},
			},
		},
	},
}

func TestGetMapsValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getMapsValueTests {
		r, e := m.GetMaps(test.path)
		if test.err {
			assert.NotNil(t, e, fmt.Sprintf("Error has been returned on %s", test.path))
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded "+test.path)
		} else {
			assert.Nil(t, e, fmt.Sprintf("NO error returned on %s (%+v)", test.path, r))
		}
		if test.expected == nil {
			assert.Nil(t, r, fmt.Sprintf("Expected nil returned on %s", test.path))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("Expected value returned on %s", test.path))
		}
	}
}

func TestGetMapsValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := []map[string]interface{}{
		map[string]interface{}{
			"foo": "bar",
		},
		map[string]interface{}{
			"bar": "baz",
		},
	}
	r, e := m.GetMaps("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path (ints)")
	assert.Equal(t, r, f, "Fallback is returned (ints)")
}

func TestGetMapsSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for i, test := range getMapsValueTests {
		r := m.GetMapsV(test.path)
		if test.err {
			assert.Nil(t, r, fmt.Sprintf("[%d] Nil result returned", i))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("[%d] Expected value returned (ACTUAL: %+v)", i, r))
		}
	}
}

/*
 * -------
 * Get: Subs (list)
 * -------
 */

var getSubsValueTests = []struct {
	path     string
	err      bool
	expected interface{}
}{
	// from single int
	{
		path:     "scalar/realint",
		err:      true,
		expected: nil,
	},
	// from single bool
	{
		path:     "bool/yes",
		err:      true,
		expected: nil,
	},
	// from actual float
	{
		path:     "scalar/realfloat",
		err:      true,
		expected: nil,
	},
	// from parsable int string
	{
		path:     "scalar/stringint",
		err:      true,
		expected: nil,
	},
	// from parsable float string
	{
		path:     "scalar/stringfloat",
		err:      true,
		expected: nil,
	},
	// from not regular string
	{
		path:     "foo/bar",
		err:      true,
		expected: nil,
	},
	// from parsable struct
	{
		path:     "foo/baz",
		err:      true,
		expected: nil,
	},
	// from array of ints
	{
		path:     "array/realints",
		err:      true,
		expected: nil,
	},
	// from array of bools
	{
		path:     "array/realbools",
		err:      true,
		expected: nil,
	},
	// from array of floats
	{
		path:     "array/realfloats",
		err:      true,
		expected: nil,
	},
	// from array of ints
	{
		path:     "array/stringints",
		err:      true,
		expected: nil,
	},
	// from array of floats
	{
		path:     "array/stringfloats",
		err:      true,
		expected: nil,
	},
	// from array of floats
	{
		path:     "array/strings",
		err:      true,
		expected: nil,
	},
	// from empty array
	{
		path:     "array/empty",
		err:      false,
		expected: []*MapPath{},
	},
	// from invalid path
	{
		path:     "x/y/z",
		err:      true,
		expected: nil,
	},
	// from invalid path
	{
		path: "mixed/array2",
		err:  false,
		expected: []*MapPath{
			NewMapPath(map[string]interface{}{
				"foo": []int{1, 2, 3, 4},
				"bar": []string{"one", "two"},
			}),
			NewMapPath(map[string]interface{}{
				"foo": []int{11, 12, 13, 14},
				"bar": []string{"five", "six"},
			}),
		},
	},
}

func TestGetSubsValue(t *testing.T) {
	m := NewMapPath(defaultTest)
	for _, test := range getSubsValueTests {
		r, e := m.GetSubs(test.path)
		if test.err {
			assert.NotNil(t, e, fmt.Sprintf("Error has been returned on %s", test.path))
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded "+test.path)
		} else {
			assert.Nil(t, e, fmt.Sprintf("NO error returned on %s (%+v)", test.path, r))
		}
		if test.expected == nil {
			assert.Nil(t, r, fmt.Sprintf("Expected nil returned on %s", test.path))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("Expected value returned on %s", test.path))
		}
	}
}

func TestGetSubsValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := []*MapPath{
		NewMapPath(map[string]interface{}{
			"foo": "bar",
		}),
		NewMapPath(map[string]interface{}{
			"bar": "baz",
		}),
	}
	r, e := m.GetSubs("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path (ints)")
	assert.Equal(t, r, f, "Fallback is returned (ints)")
}

func TestGetSubsSingleContext(t *testing.T) {
	m := NewMapPath(defaultTest)
	for i, test := range getSubsValueTests {
		r := m.GetSubsV(test.path)
		if test.err {
			assert.Nil(t, r, fmt.Sprintf("[%d] Nil result returned", i))
		} else {
			assert.Equal(t, test.expected, r, fmt.Sprintf("[%d] Expected value returned (ACTUAL: %+v)", i, r))
		}
	}
}

/*
 * -------
 * Get: Array
 * -------
 */

func TestGetUnsupportedArrayValueReturnsError(t *testing.T) {
	m := NewMapPath(defaultTest)
	r, ok, e := m.GetArray(reflect.TypeOf(byte(0)), "array/realints")
	_, isaUnsupportedTypeError := e.(UnsupportedTypeError)
	assert.NotNil(t, e, "Error returned on unsupported type")
	assert.True(t, isaUnsupportedTypeError, "Unsupported type error has been returned")
	assert.Nil(t, r, "Result is nil")
	assert.False(t, ok, "Not been found")
}

/*
 * -------
 * Error
 * -------
 */

var errorNotFoundErrorFormatTests = [][]string{
	[]string{"foo/bar/baz", "The path \"foo/bar/baz\" does not exist"},
	[]string{"f", "The path \"f\" does not exist"},
}

func TestNotFoundErrorFormat(t *testing.T) {
	for _, test := range errorNotFoundErrorFormatTests {
		err := NotFoundError(test[0])
		assert.Equal(t, err.Error(), test[1], "Error correctly formatted")
	}
}

var invalidTypeErrorFormatTests = []struct {
	val    interface{}
	expect string
	msg    string
}{
	{
		val:    123,
		expect: "foo",
		msg:    "Could not cast int into foo",
	},
	{
		val:    0.0,
		expect: "foo",
		msg:    "Could not cast float64 into foo",
	},
	{
		val:    []int{1, 2},
		expect: "foo",
		msg:    "Could not cast []int into foo",
	},
	{
		val:    map[string]interface{}{"x": "y"},
		expect: "foo",
		msg:    "Could not cast map[string]interface {} into foo",
	},
}

func TestInvalidTypeErrorFormat(t *testing.T) {
	for _, test := range invalidTypeErrorFormatTests {
		err := &InvalidTypeError{test.val, test.expect}
		assert.Equal(t, err.Error(), test.msg, "Error correctly formatted")
	}
}

var errorUnsupportedTypeTests = [][]string{
	[]string{"int8", "Type int8 is not supported"},
	[]string{"foo-bar", "Type foo-bar is not supported"},
}

func TestUnsupportedType(t *testing.T) {
	for _, test := range errorUnsupportedTypeTests {
		err := UnsupportedTypeError(test[0])
		assert.Equal(t, err.Error(), test[1], "Error correctly formatted")
	}
}

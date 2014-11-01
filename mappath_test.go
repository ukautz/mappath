package mappath

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var defaultTest = map[string]interface{}{
	"hello": "world",
	"foo": map[string]interface{}{
		"bar": "baz",
		"baz": map[string]interface{}{
			"bam": 42,
		},
	},
	"array": []int{1, 2, 3, 4},
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
	"some": map[string]interface{}{
		"stringint":   "123",
		"stringfloat": "123.456",
		"realint":     123,
		"realfloat":   123.456,
	},
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
		path:   "array/0",
		expect: 1,
		from:   defaultTest,
	},
	{
		path:   "array/3",
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
		path:   "array",
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
		path:     "some/realint",
		err:      false,
		expected: 123,
	},
	// from actual float
	{
		path:     "some/realfloat",
		err:      false,
		expected: 123,
	},
	// from parsable int string
	{
		path:     "some/stringint",
		err:      false,
		expected: 123,
	},
	// from parsable float string
	{
		path:     "some/stringfloat",
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
		path:     "array",
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
		assert.Equal(t, r, test.expected, "Expected value returned")
	}
}

func TestGetIntValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := 345
	r, e := m.GetInt("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path")
	assert.Equal(t, r, f, "Fallback is returned")
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
		path:     "some/realint",
		err:      false,
		expected: 123.0,
	},
	// from actual float
	{
		path:     "some/realfloat",
		err:      false,
		expected: 123.456,
	},
	// from parsable int string
	{
		path:     "some/stringint",
		err:      false,
		expected: 123.0,
	},
	// from parsable float string
	{
		path:     "some/stringfloat",
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
		path:     "array",
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
		assert.Equal(t, r, test.expected, "Expected value returned")
	}
}

func TestGetFloatValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := 345.678
	r, e := m.GetFloat("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path")
	assert.Equal(t, r, f, "Fallback is returned")
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
		path:     "some/realint",
		err:      false,
		expected: "123",
	},
	// from actual float
	{
		path:     "some/realfloat",
		err:      false,
		expected: "123.456000000",
	},
	// from parsable int string
	{
		path:     "some/stringint",
		err:      false,
		expected: "123",
	},
	// from parsable float string
	{
		path:     "some/stringfloat",
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
		path:     "array",
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
		assert.Equal(t, r, test.expected, "Expected value returned "+test.path)
	}
}

func TestGetStringValueFallback(t *testing.T) {
	m := NewMapPath(map[string]interface{}{})
	f := "foo"
	r, e := m.GetString("x/y/z", f)
	assert.Nil(t, e, "No error when fallback used on invalid path")
	assert.Equal(t, r, f, "Fallback is returned")
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
		path:     "some/realint",
		err:      true,
		expected: nil,
	},
	// from actual float
	{
		path:     "some/realfloat",
		err:      true,
		expected: nil,
	},
	// from parsable int string
	{
		path:     "some/stringint",
		err:      true,
		expected: nil,
	},
	// from parsable float string
	{
		path:     "some/stringfloat",
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
		path:     "array",
		err:      true,
		expected: nil,
	},
	// from invalid path
	{
		path:     "x/y/z",
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
			assert.NotNil(t, e, "Error returned OK "+test.path)
			assert.IsType(t, reflect.TypeOf(&InvalidTypeError{}), reflect.TypeOf(e), "Correct error responded "+test.path)
		} else {
			assert.Nil(t, e, "NO error returned "+test.path)
		}
		if test.expected == nil {
			assert.Nil(t, r, "Expected nil returned "+test.path)
		} else {
			assert.Equal(t, r, test.expected, "Expected value returned "+test.path)
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

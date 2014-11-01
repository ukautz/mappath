package mappath

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromValidJson(t *testing.T) {
	r, e := FromJsonFile("resources/ok.json")
	assert.Nil(t, e, "No error returned")
	d, e := r.GetString("foo")
	assert.Nil(t, e, "foo key found")
	assert.Equal(t, "bar", d, "bar value returned")
}

func TestFromInvalidJsonFile(t *testing.T) {
	r, e := FromJsonFile("resources/invalid.json")
	assert.NotNil(t, e, "Error has been returned")
	assert.Nil(t, r, "No result is returned")
}

func TestFromUnsupportedButValidJsonFile(t *testing.T) {
	r, e := FromJsonFile("resources/fail.json")
	assert.NotNil(t, e, "Error has been returned")
	assert.Nil(t, r, "No result is returned")
}

func TestFromMissingJsonFile(t *testing.T) {
	r, e := FromJsonFile("resources/missing.json")
	assert.NotNil(t, e, "Error has been returned")
	assert.Nil(t, r, "No result is returned")
}

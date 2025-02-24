package hw09structvalidator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLen(t *testing.T) {
	valid, err := checkLen("test", "4")
	assert.NoError(t, err)
	assert.True(t, valid)

	valid, err = checkLen("test", "5")
	assert.NoError(t, err)
	assert.False(t, valid)

	valid, err = checkLen("test", "abc")
	assert.Error(t, err)
	assert.False(t, valid)

	valid, err = checkLen("test", "-1")
	assert.Error(t, err)
	assert.False(t, valid)
	assert.Contains(t, err.Error(), "negative len")
}

func TestCheckRegexp(t *testing.T) {
	valid, err := checkRegexp("hello", "h.*o")
	assert.NoError(t, err)
	assert.True(t, valid)

	valid, err = checkRegexp("hello", "^h.*z$")
	assert.NoError(t, err)
	assert.False(t, valid)

	valid, err = checkRegexp("hello", "(*")
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestCheckIn(t *testing.T) {
	valid, err := checkIn(reflect.ValueOf("a"), "a,b,c")
	assert.NoError(t, err)
	assert.True(t, valid)

	valid, err = checkIn(reflect.ValueOf("d"), "a,b,c")
	assert.NoError(t, err)
	assert.False(t, valid)

	valid, err = checkIn(reflect.ValueOf(3), "1,2,3")
	assert.NoError(t, err)
	assert.True(t, valid)

	valid, err = checkIn(reflect.ValueOf(4), "1,2,3")
	assert.NoError(t, err)
	assert.False(t, valid)

	valid, err = checkIn(reflect.ValueOf(3.14), "3.14,2.71")
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestCheckMin(t *testing.T) {
	valid, err := checkMin(5, "3")
	assert.NoError(t, err)
	assert.True(t, valid)

	valid, err = checkMin(2, "3")
	assert.NoError(t, err)
	assert.False(t, valid)

	valid, err = checkMin(2, "abc")
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestCheckMax(t *testing.T) {

	valid, err := checkMax(3, "5")
	assert.NoError(t, err)
	assert.True(t, valid)

	valid, err = checkMax(7, "5")
	assert.NoError(t, err)
	assert.False(t, valid)

	valid, err = checkMax(7, "abc")
	assert.Error(t, err)
	assert.False(t, valid)
}

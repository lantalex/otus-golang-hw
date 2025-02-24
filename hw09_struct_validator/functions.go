package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func checkLen(value string, args string) (valid bool, err error) {
	expected, err := strconv.Atoi(args)
	if err != nil {
		return false, err
	}

	if expected < 0 {
		return false, fmt.Errorf("negative len: %d", expected)
	}

	if len(value) != expected {
		return false, nil
	}

	return true, nil
}

func checkRegexp(value string, args string) (valid bool, err error) {

	r, err := regexp.Compile(args)
	if err != nil {
		return false, err
	}

	return r.MatchString(value), nil
}

func checkIn(value reflect.Value, args string) (valid bool, err error) {
	var str string

	switch value.Kind() {
	case reflect.String:
		str = value.String()
	case reflect.Int:
		str = strconv.FormatInt(value.Int(), 10)
	default:
		return false, fmt.Errorf("is not a string/int")
	}

	for _, item := range strings.Split(args, ",") {
		if item == str {
			return true, nil
		}
	}

	return false, nil
}

func checkMin(value int, args string) (valid bool, err error) {
	expected, err := strconv.Atoi(args)
	if err != nil {
		return false, err
	}

	if value < expected {
		return false, nil
	}

	return true, nil
}

func checkMax(value int, args string) (valid bool, err error) {
	expected, err := strconv.Atoi(args)
	if err != nil {
		return false, err
	}

	if value > expected {
		return false, nil
	}

	return true, nil
}

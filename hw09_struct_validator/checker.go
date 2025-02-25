package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type Checker interface {
	Check(v reflect.Value, args string) (valid bool, err error)
}

type ReflectChecker struct {
	function func(reflect.Value, string) (bool, error)
}

func (c *ReflectChecker) Check(v reflect.Value, args string) (bool, error) {
	return c.function(v, args)
}

type FunctionChecker[T string | int | float64] struct {
	function func(T, string) (bool, error)
}

func (c *FunctionChecker[T]) Check(v reflect.Value, args string) (valid bool, err error) {
	u, ok := v.Interface().(T)

	if !ok {
		var t T
		return false, fmt.Errorf("unsupported type: got %T, expected %T", v.Interface(), t)
	}

	return c.function(u, args)
}

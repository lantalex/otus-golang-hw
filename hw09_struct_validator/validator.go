package hw09structvalidator

import (
	"bytes"
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	buffer := bytes.Buffer{}
	for _, val := range v {
		buffer.WriteString(val.Field)
		buffer.WriteString(": ")
		buffer.WriteString(val.Err.Error())
		buffer.WriteString("\n")
	}

	return fmt.Sprint(buffer.String())
}

func Validate(v interface{}) error {
	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("value must be a struct, %T instead", t)
	}

	engine := NewEmptyEngine()

	engine.addChecker("notNan", &FunctionChecker[float64]{checkNotNan})
	engine.addChecker("len", &FunctionChecker[string]{checkLen})
	engine.addChecker("regexp", &FunctionChecker[string]{checkRegexp})
	engine.addChecker("min", &FunctionChecker[int]{checkMin})
	engine.addChecker("max", &FunctionChecker[int]{checkMax})
	engine.addChecker("in", &ReflectChecker{checkIn})

	engine.addNestedCheck("nested")

	ve, err := engine.Validate(t.Name(), "validate:\"nested\"", reflect.ValueOf(v))
	if err != nil {
		return err
	}

	return ve
}

package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Engine struct {
	checkers map[string]Checker
}

func NewEmptyEngine() *Engine {
	return &Engine{checkers: make(map[string]Checker)}
}

func (engine *Engine) Validate(path string, tag reflect.StructTag, v reflect.Value) (ValidationErrors, error) {

	checks, err := engine.resolveChecks(path, tag)

	if err != nil {
		return nil, err
	}

	var errs ValidationErrors

	for _, check := range checks {

		ve, err := engine.invoke(check, v)

		if err != nil {
			return nil, err
		}

		errs = append(errs, ve...)
	}
	return errs, nil
}

func (engine *Engine) resolveChecks(path string, tag reflect.StructTag) ([]engineInternalCheck, error) {

	validationTag, ok := tag.Lookup("Validate")
	if !ok || validationTag == "" {
		return nil, nil
	}

	var checks []engineInternalCheck

	for _, item := range strings.Split(validationTag, "|") {
		name, args, _ := strings.Cut(item, ":")

		rule, ok := engine.checkers[name]

		if !ok {
			return nil, fmt.Errorf("%s: unknown check: %s", path, name)
		}

		checks = append(checks, engineInternalCheck{rule, name, args, path, tag})
	}

	return checks, nil
}

func (engine *Engine) invoke(ec engineInternalCheck, v reflect.Value) (ValidationErrors, error) {

	var errs ValidationErrors

	path := ec.fieldPath
	tag := ec.fieldTag

	defer func() {
		ec.fieldPath = path
		ec.fieldTag = tag
	}()

	switch v.Kind() {
	case reflect.String, reflect.Int:
		return ec.check(v)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {

			ec.fieldPath = path + "[" + strconv.Itoa(i) + "]"
			ve, err := ec.check(v.Index(i))

			if err != nil {
				return nil, err
			}

			errs = append(errs, ve...)
		}
		return errs, nil
	case reflect.Struct:

		for i := 0; i < v.NumField(); i++ {

			ec.fieldPath = path + "." + v.Type().Field(i).Name
			ec.fieldTag = v.Type().Field(i).Tag

			ve, err := ec.check(v.Field(i))

			if err != nil {
				return nil, err
			}

			errs = append(errs, ve...)
		}
		return errs, nil
	default:
		return nil, fmt.Errorf("%s: unsupported type: %s", ec.fieldPath, v.Type())
	}
}

type EngineChecker struct {
	engine *Engine
	path   string
	tag    reflect.StructTag
}

func (c *EngineChecker) Check(v reflect.Value, _ string) (valid bool, err error) {

	ve, err := c.engine.Validate(c.path, c.tag, v)

	if err != nil {
		return false, err
	}

	if ve != nil {
		return false, ve
	}

	return true, nil
}

type engineInternalCheck struct {
	checker   Checker
	checkName string
	checkArgs string

	fieldPath string
	fieldTag  reflect.StructTag
}

func (ec *engineInternalCheck) check(v reflect.Value) (ValidationErrors, error) {

	ok, err := ec.checker.Check(v, ec.checkArgs)

	if err != nil {
		var ve ValidationErrors
		if errors.As(err, &ve) {
			return ve, nil
		} else {
			return nil, fmt.Errorf(
				"%s: fatal error during \"%s:%s\" for value \"%v\": %w",
				ec.fieldPath,
				ec.checkName,
				ec.checkArgs,
				v,
				err,
			)
		}
	}

	if !ok {
		return ValidationErrors{ValidationError{
			Field: ec.fieldPath,
			Err: fmt.Errorf(
				"\"%v\" does not satisfy \"%s:%s\"",
				v,
				ec.checkName,
				ec.checkArgs,
			),
		}}, nil
	}

	return nil, nil
}

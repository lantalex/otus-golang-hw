package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Engine struct {
	checkers map[string]InternalChecker
}

func NewEmptyEngine() *Engine {
	return &Engine{checkers: make(map[string]InternalChecker)}
}

func (engine *Engine) addChecker(name string, checker Checker) {
	engine.checkers[name] = &ProxyChecker{checker: checker, name: name}
}

func (engine *Engine) addNestedCheck(name string) {
	engine.checkers[name] = &EngineChecker{engine: engine}
}

func (engine *Engine) Validate(path string, tag reflect.StructTag, v reflect.Value) (ValidationErrors, error) {
	checks, allArgs, err := engine.resolveCheckersAndArgs(path, tag)
	if err != nil {
		return nil, err
	}

	var errs ValidationErrors

	for i, check := range checks {
		ve, err := engine.invoke(path, tag, check, allArgs[i], v)
		if err != nil {
			return nil, err
		}

		errs = append(errs, ve...)
	}
	return errs, nil
}

func (engine *Engine) resolveCheckersAndArgs(path string, tag reflect.StructTag) ([]InternalChecker, []string, error) {
	validationTag, ok := tag.Lookup("validate")
	if !ok || validationTag == "" {
		return nil, nil, nil
	}

	checks := make([]InternalChecker, 0)
	allArgs := make([]string, 0)

	for _, item := range strings.Split(validationTag, "|") {
		name, args, _ := strings.Cut(item, ":")

		internalCheck, ok := engine.checkers[name]

		if !ok {
			return nil, nil, fmt.Errorf("%s: unknown check: %s", path, name)
		}

		checks = append(checks, internalCheck)
		allArgs = append(allArgs, args)
	}

	return checks, allArgs, nil
}

func (engine *Engine) invoke(
	path string,
	tag reflect.StructTag,
	c InternalChecker,
	args string,
	v reflect.Value,
) (ValidationErrors, error) {
	var errs ValidationErrors

	if v.Kind() == reflect.String || v.Kind() == reflect.Int || v.Kind() == reflect.Float64 {
		return c.check(v, args, path, tag)
	}

	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			ve, err := c.check(v.Index(i), args, path+"["+strconv.Itoa(i)+"]", tag)
			if err != nil {
				return nil, err
			}
			errs = append(errs, ve...)
		}
		return errs, nil
	}

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			ve, err := c.check(v.Field(i), args, path+"."+v.Type().Field(i).Name, v.Type().Field(i).Tag)
			if err != nil {
				return nil, err
			}

			errs = append(errs, ve...)
		}
		return errs, nil
	}
	return nil, fmt.Errorf("%s: unsupported type: %s", path, v.Type())
}

type InternalChecker interface {
	check(v reflect.Value, args string, path string, tag reflect.StructTag) (ValidationErrors, error)
}

type ProxyChecker struct {
	checker Checker
	name    string
}

func (ec *ProxyChecker) check(
	v reflect.Value,
	args string,
	path string,
	_ reflect.StructTag,
) (ValidationErrors, error) {
	ok, err := ec.checker.Check(v, args)
	if err != nil {
		var ve ValidationErrors
		if errors.As(err, &ve) {
			return ve, nil
		}
		return nil, fmt.Errorf(
			"%s: fatal error during \"%s:%s\" for value \"%v\": %w",
			path,
			ec.name,
			args,
			v,
			err,
		)
	}

	if !ok {
		return ValidationErrors{ValidationError{
			Field: path,
			Err: fmt.Errorf(
				"\"%v\" does not satisfy \"%s:%s\"",
				v,
				ec.name,
				args,
			),
		}}, nil
	}

	return nil, nil
}

type EngineChecker struct {
	engine *Engine
}

func (c *EngineChecker) check(v reflect.Value, _ string, path string, tag reflect.StructTag) (ValidationErrors, error) {
	return c.engine.Validate(path, tag, v)
}

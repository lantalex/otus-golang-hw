package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert" //nolint:depguard
	"github.com/stretchr/testify/mock"   //nolint:depguard
)

type MockChecker struct {
	mock.Mock
}

func (m *MockChecker) Check(v reflect.Value, args string) (valid bool, err error) {
	argsMock := m.Called(v, args)
	return argsMock.Bool(0), argsMock.Error(1)
}

func TestEngine_Validate(t *testing.T) {
	tests := []struct {
		name                   string
		path                   string
		tag                    reflect.StructTag
		value                  reflect.Value
		mockSetup              func(*MockChecker, reflect.Value)
		expectedValidationErrs ValidationErrors
		expectedErr            error
	}{
		{
			name:  "string successful validation",
			path:  "field",
			tag:   reflect.StructTag(`Validate:"mock:arg"`),
			value: reflect.ValueOf("test"),
			mockSetup: func(m *MockChecker, value reflect.Value) {
				m.On("Check", value, "arg").Return(true, nil).Once()
			},
			expectedValidationErrs: nil,
			expectedErr:            nil,
		},
		{
			name:  "slice successful validation",
			path:  "field",
			tag:   reflect.StructTag(`Validate:"mock:arg"`),
			value: reflect.ValueOf([]string{"val1", "val2"}),
			mockSetup: func(m *MockChecker, value reflect.Value) {
				m.On("Check", value.Index(0), "arg").Return(true, nil).Once()
				m.On("Check", value.Index(1), "arg").Return(true, nil).Once()
			},
			expectedValidationErrs: nil,
			expectedErr:            nil,
		},
		{
			name: "struct successful validation",
			path: "field",
			tag:  reflect.StructTag(`Validate:"mock:arg"`),
			value: reflect.ValueOf(struct {
				arg1 string
				arg2 string
			}{"arg1val", "arg2val"}),
			mockSetup: func(m *MockChecker, value reflect.Value) {
				m.On("Check", value.Field(0), "arg").Return(true, nil).Once()
				m.On("Check", value.Field(1), "arg").Return(true, nil).Once()
			},
			expectedValidationErrs: nil,
			expectedErr:            nil,
		},
		{
			name:  "string unsuccessful validation",
			path:  "field",
			tag:   reflect.StructTag(`Validate:"mock:arg"`),
			value: reflect.ValueOf("test"),
			mockSetup: func(m *MockChecker, value reflect.Value) {
				m.On("Check", value, "arg").Return(false, nil).Once()
			},
			expectedValidationErrs: ValidationErrors{ValidationError{Field: "field", Err: errors.New("\"test\" does not satisfy \"mock:arg\"")}},
			expectedErr:            nil,
		},
		{
			name:  "slice unsuccessful validation",
			path:  "slice",
			tag:   reflect.StructTag(`Validate:"mock:arg"`),
			value: reflect.ValueOf([]string{"val1", "val2"}),
			mockSetup: func(m *MockChecker, value reflect.Value) {
				m.On("Check", value.Index(0), "arg").Return(true, nil).Once()
				m.On("Check", value.Index(1), "arg").Return(false, nil).Once()
			},
			expectedValidationErrs: ValidationErrors{ValidationError{Field: "slice[1]", Err: errors.New("\"val2\" does not satisfy \"mock:arg\"")}},
			expectedErr:            nil,
		},
		{
			name: "struct unsuccessful validation",
			path: "struct",
			tag:  reflect.StructTag(`Validate:"mock:arg"`),
			value: reflect.ValueOf(struct {
				field1 string
				field2 string
			}{"field1val", "field2val"}),
			mockSetup: func(m *MockChecker, value reflect.Value) {
				m.On("Check", value.Field(0), "arg").Return(true, nil).Once()
				m.On("Check", value.Field(1), "arg").Return(false, nil).Once()
			},
			expectedValidationErrs: ValidationErrors{ValidationError{Field: "struct.field2", Err: errors.New("\"field2val\" does not satisfy \"mock:arg\"")}},
			expectedErr:            nil,
		},
		{
			name:  "validation error",
			path:  "field",
			tag:   reflect.StructTag(`Validate:"mock:arg"`),
			value: reflect.ValueOf("test"),
			mockSetup: func(m *MockChecker, value reflect.Value) {
				m.On("Check", value, "arg").Return(false, errors.New("validation error")).Once()
			},
			expectedValidationErrs: nil,
			expectedErr:            errors.New("field: fatal error during \"mock:arg\" for value \"test\": validation error"),
		},

		{
			name:  "unknown check",
			path:  "field",
			tag:   reflect.StructTag(`Validate:"unknown"`),
			value: reflect.ValueOf("test"),
			mockSetup: func(m *MockChecker, value reflect.Value) {
				m.AssertNotCalled(t, "Check", value, "arg")
			},
			expectedValidationErrs: nil,
			expectedErr:            errors.New("field: unknown check: unknown"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			engine := NewEmptyEngine()

			mockChecker := new(MockChecker)
			engine.checkers["mock"] = mockChecker

			if tt.mockSetup != nil {
				tt.mockSetup(mockChecker, tt.value)
			}

			ve, err := engine.Validate(tt.path, tt.tag, tt.value)

			if tt.expectedValidationErrs != nil {
				assert.NotNil(t, ve)
				assert.Equal(t, tt.expectedValidationErrs, ve)
			} else {
				assert.Nil(t, ve)
			}

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockChecker.AssertExpectations(t)
		})
	}
}

func TestEngine_resolveChecks(t *testing.T) {

	engine := NewEmptyEngine()

	mockChecker := new(MockChecker)
	engine.checkers["mock"] = mockChecker

	tests := []struct {
		name        string
		path        string
		tag         reflect.StructTag
		expected    []engineInternalCheck
		expectedErr error
	}{
		{
			name: "valid tag",
			path: "field",
			tag:  reflect.StructTag(`Validate:"mock:arg"`),
			expected: []engineInternalCheck{
				{
					checker:   mockChecker,
					checkName: "mock",
					checkArgs: "arg",
					fieldPath: "field",
					fieldTag:  reflect.StructTag(`Validate:"mock:arg"`),
				},
			},
			expectedErr: nil,
		},
		{
			name:        "empty tag",
			path:        "field",
			tag:         reflect.StructTag(``),
			expected:    nil,
			expectedErr: nil,
		},
		{
			name:        "unknown check",
			path:        "field",
			tag:         reflect.StructTag(`Validate:"unknown:arg"`),
			expected:    nil,
			expectedErr: fmt.Errorf("field: unknown check: unknown"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			checks, err := engine.resolveChecks(tt.path, tt.tag)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, checks)
			}
		})
	}
}

func TestEngine_invoke(t *testing.T) {

	tests := []struct {
		name                   string
		ec                     engineInternalCheck
		value                  reflect.Value
		mockSetup              func(*MockChecker)
		expectedValidationErrs ValidationErrors
		expectedErr            error
	}{
		{
			name:  "successful validation",
			value: reflect.ValueOf("test"),
			mockSetup: func(m *MockChecker) {
				m.On("Check", reflect.ValueOf("test"), "arg").Return(true, nil)
			},
			expectedValidationErrs: nil,
			expectedErr:            nil,
		},
		{
			name:  "unsuccessful validation",
			value: reflect.ValueOf("test"),
			mockSetup: func(m *MockChecker) {
				m.On("Check", reflect.ValueOf("test"), "arg").Return(false, nil)
			},
			expectedValidationErrs: ValidationErrors{ValidationError{Field: "field", Err: errors.New("\"test\" does not satisfy \"mock:arg\"")}},
			expectedErr:            nil,
		},
		{
			name:  "validation error",
			value: reflect.ValueOf("test"),
			mockSetup: func(m *MockChecker) {
				m.On("Check", reflect.ValueOf("test"), "arg").Return(false, errors.New("validation error"))
			},
			expectedValidationErrs: nil,
			expectedErr:            errors.New("field: fatal error during \"mock:arg\" for value \"test\": validation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockChecker := new(MockChecker)

			if tt.mockSetup != nil {
				tt.mockSetup(mockChecker)
			}

			engine := NewEmptyEngine()
			engine.checkers["mock"] = mockChecker

			ec := engineInternalCheck{
				checker:   mockChecker,
				checkName: "mock",
				checkArgs: "arg",
				fieldPath: "field",
				fieldTag:  reflect.StructTag(`Validate:"mock:arg"`),
			}

			ve, err := engine.invoke(ec, tt.value)

			if tt.expectedValidationErrs != nil {
				assert.NotNil(t, ve)
				assert.Equal(t, tt.expectedValidationErrs, ve)
				assert.Nil(t, err)
			} else {
				assert.Nil(t, ve)
			}

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, ve)
			} else {
				assert.NoError(t, err)
			}

			mockChecker.AssertExpectations(t)
		})
	}
}

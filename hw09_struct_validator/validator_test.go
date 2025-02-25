package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert" //nolint:depguard
)

// Test the function on different structures and other types.
type (
	UserRole string

	NestedStruct struct {
		IntField    int    `validate:"min:1"`
		StringField string `validate:"len:1"`
	}

	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int          `validate:"min:18|max:50"`
		Email  string       `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole     `validate:"in:admin,stuff"`
		Phones []string     `validate:"len:11"`
		Nested NestedStruct `validate:"nested"`
		Score  float64      `validate:"notNan"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "this is valid   id  with length = 36",
				Name:   "John Smith",
				Age:    27,
				Email:  "valid@email.com",
				Role:   UserRole("stuff"),
				Phones: []string{"12345678901"},
				Nested: NestedStruct{IntField: 1, StringField: "1"},
				Score:  1.0,
				meta:   json.RawMessage(`{"key":"value"}`),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "bad id",
				Name:   "Alice",
				Age:    7,
				Email:  "invalidemail",
				Role:   UserRole("child"),
				Phones: []string{"12345678901", "0"},
				Nested: NestedStruct{IntField: -1, StringField: "42"},
				Score:  math.NaN(),
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "User.ID",
					Err:   errors.New("\"bad id\" does not satisfy \"len:36\""),
				},
				ValidationError{
					Field: "User.Age",
					Err:   errors.New("\"7\" does not satisfy \"min:18\""),
				},
				ValidationError{
					Field: "User.Email",
					Err:   errors.New("\"invalidemail\" does not satisfy \"regexp:^\\w+@\\w+\\.\\w+$\""),
				},
				ValidationError{
					Field: "User.Role",
					Err:   errors.New("\"child\" does not satisfy \"in:admin,stuff\""),
				},
				ValidationError{
					Field: "User.Phones[1]",
					Err:   errors.New("\"0\" does not satisfy \"len:11\""),
				},
				ValidationError{
					Field: "User.Nested.IntField",
					Err:   errors.New("\"-1\" does not satisfy \"min:1\""),
				},
				ValidationError{
					Field: "User.Nested.StringField",
					Err:   errors.New("\"42\" does not satisfy \"len:1\""),
				},
				ValidationError{
					Field: "User.Score",
					Err:   errors.New("\"NaN\" does not satisfy \"notNan:\""),
				},
			},
		},
		{
			in:          App{Version: "1.0.0"},
			expectedErr: nil,
		},
		{
			in: App{Version: "1.0"},
			expectedErr: ValidationErrors{ValidationError{
				Field: "App.Version",
				Err:   errors.New("\"1.0\" does not satisfy \"len:5\""),
			}},
		},
		{
			in:          Token{Header: []byte(`{"key":"value"}`)},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 200, Body: "ok"},
			expectedErr: nil,
		},
		{
			in: Response{Code: 418, Body: "I'm a teapot"},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Response.Code",
				Err:   errors.New("\"418\" does not satisfy \"in:200,404,500\""),
			}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expectedErr, err)
			}
		})
	}
}

func TestFatalError(t *testing.T) {
	type Bad struct {
		Str string `validate:"len:foo"`
	}

	err := Validate(Bad{Str: "bar"})

	assert.NotNil(t, err)
	assert.Equal(t,
		errors.New("Bad.Str: fatal error during \"len:foo\" for value \"bar\": "+
			"strconv.Atoi: parsing \"foo\": invalid syntax").Error(),
		err.Error())
}

package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: `a4bc2d5e`, expected: `aaaabccddddde`},
		{input: `abccd`, expected: `abccd`},
		{input: ``, expected: ""},
		{input: `aaa0b`, expected: `aab`},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `юникод3:0 \﷽3`, expected: `юникоддд ﷽﷽﷽`},
		{input: `флаг: 🇷🇺3 и 🇸🇭0`, expected: `флаг: 🇷🇺🇷🇺🇷🇺 и `},
		{input: `한5 한1 한2 한 3 글4 한\0`, expected: `한한한한한 한 한한 한    글글글글 한0`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{
		`3abc`,
		`45`,
		`aaa10b`,
		`asda\`,
		`\`,
		"строка в некорректной UTF-8 кодировке:\x80",
	}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

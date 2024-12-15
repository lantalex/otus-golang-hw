package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/rivo/uniseg" //nolint:depguard
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if !utf8.ValidString(str) {
		return "", ErrInvalidString
	}

	var res strings.Builder
	chunk := ""
	backslash := false

	gr := uniseg.NewGraphemes(str)
	for gr.Next() {
		beg, end := gr.Positions()
		next := str[beg:end]

		// backslash case
		if next == `\` {
			if backslash {
				chunk = str[beg:end]
				backslash = false
			} else {
				res.WriteString(chunk)
				chunk = ""
				backslash = true
			}
			continue
		}

		// digit case
		if digitValue, err := strconv.Atoi(next); err == nil {
			if backslash {
				res.WriteString(chunk)
				chunk = str[beg:end]
				backslash = false
			} else {
				if chunk == "" {
					return "", ErrInvalidString
				}

				res.WriteString(strings.Repeat(chunk, digitValue))
				chunk = ""
			}
			continue
		}

		// regular symbol case
		if backslash {
			return "", ErrInvalidString
		}
		res.WriteString(chunk)
		chunk = str[beg:end]
		backslash = false
	}

	if backslash {
		return "", ErrInvalidString
	}

	res.WriteString(chunk)

	return res.String(), nil
}

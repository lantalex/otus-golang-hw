package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/rivo/uniseg"
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

		if digitValue, err := strconv.Atoi(next); err == nil && !backslash {
			if chunk == "" {
				return "", ErrInvalidString
			}

			res.WriteString(strings.Repeat(chunk, digitValue))
			chunk = ""
			backslash = false
			continue
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

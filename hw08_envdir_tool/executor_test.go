package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestRunCmd(t *testing.T) {
	t.Run("echo", func(t *testing.T) {
		env := Environment{
			"BAR":   {"bar", false},
			"UNSET": {"", true},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
		}

		code, err := RunCmd([]string{"./testdata/echo.sh"}, env)
		require.NoError(t, err)
		require.Equal(t, CmdOk, code)
	})

	t.Run("use code return from cmd - false", func(t *testing.T) {
		code, err := RunCmd([]string{"false"}, Environment{})
		require.NoError(t, err)
		require.Equal(t, 1, code)
	})

	t.Run("command not found", func(t *testing.T) {
		code, err := RunCmd([]string{"bad_command"}, Environment{})
		require.Error(t, err)
		require.Equal(t, InternalError, code)
	})
}

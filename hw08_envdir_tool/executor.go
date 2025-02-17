package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const (
	CmdOk         = 0
	InternalError = 1
)

func RunCmd(cmd []string, env Environment) (returnCode int, err error) {
	for name, value := range env {
		if value.NeedRemove {
			if err := os.Unsetenv(name); err != nil {
				return InternalError, fmt.Errorf("failed to unset environment variable %s: %w", name, err)
			}
		} else {
			if err := os.Setenv(name, value.Value); err != nil {
				return InternalError, fmt.Errorf("failed to set environment variable %s: %w", name, err)
			}
		}
	}

	if len(cmd) == 0 {
		return InternalError, fmt.Errorf("no command specified")
	}

	toExecute := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	toExecute.Stdin = os.Stdin
	toExecute.Stdout = os.Stdout
	toExecute.Stderr = os.Stderr
	toExecute.Env = os.Environ()

	if err := toExecute.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), nil
		}
		return InternalError, err
	}

	return CmdOk, nil
}

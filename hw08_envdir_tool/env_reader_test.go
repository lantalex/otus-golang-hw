package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestReadDir(t *testing.T) {
	t.Run("regular case", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t,
			Environment{
				"BAR":   {"bar", false},
				"EMPTY": {"", false},
				"FOO":   {"   foo\nwith new line", false},
				"HELLO": {"\"hello\"", false},
				"UNSET": {"", true},
			}, env)
	})

	t.Run("dir does not exists", func(t *testing.T) {
		_, err := os.Stat("bad_dir")
		require.True(t, os.IsNotExist(err))
		_, err = ReadDir("bad_dir")
		require.Error(t, err)
	})

	t.Run("not a dir", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "not-a-dir")
		require.NoError(t, err)
		defer func() {
			err = os.Remove(tmpFile.Name())
			require.NoError(t, err)
		}()

		_, err = ReadDir(tmpFile.Name())
		require.Error(t, err)
	})

	t.Run("empty dir", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "empty-dir")
		require.NoError(t, err)
		defer func() {
			err = os.RemoveAll(tmpDir)
			require.NoError(t, err)
		}()

		env, err := ReadDir(tmpDir)
		require.NoError(t, err)
		require.Equal(t, Environment{}, env)
	})

	t.Run("ignore nested dirs", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "empty-dir")
		require.NoError(t, err)
		defer func() {
			err = os.RemoveAll(tmpDir)
			require.NoError(t, err)
		}()

		tmpNestedDir, err := os.MkdirTemp(tmpDir, "nested-dir")
		require.NoError(t, err)
		defer func() {
			err = os.RemoveAll(tmpNestedDir)
			require.NoError(t, err)
		}()

		tmpFile, err := os.Create(filepath.Join(tmpDir, "FOO"))
		require.NoError(t, err)
		defer func() {
			err = os.Remove(tmpFile.Name())
			require.NoError(t, err)
		}()

		_, err = tmpFile.WriteString("BAR42")
		require.NoError(t, err)

		env, err := ReadDir(tmpDir)
		require.NoError(t, err)
		require.Equal(t, Environment{"FOO": {"BAR42", false}},
			env)
	})

	t.Run("ignore bad filename", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "empty-dir")
		require.NoError(t, err)
		defer func() {
			err = os.RemoveAll(tmpDir)
			require.NoError(t, err)
		}()

		tmpFile, err := os.Create(filepath.Join(tmpDir, "FOO=ignored"))
		require.NoError(t, err)
		defer func() {
			err = os.Remove(tmpFile.Name())
			require.NoError(t, err)
		}()

		_, err = tmpFile.WriteString("ignore_me")
		require.NoError(t, err)

		env, err := ReadDir(tmpDir)
		require.NoError(t, err)
		require.Equal(t, Environment{}, env)
	})
}

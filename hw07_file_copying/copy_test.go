package main

import (
	"bytes"
	"errors"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyEntireFile(t *testing.T) {
	fromFile := filepath.Join("testdata", "input.txt")
	toFile := filepath.Join(t.TempDir(), "copy.txt")

	defer func(name string) {
		_ = os.Remove(name)
	}(toFile)

	if err := Copy(fromFile, toFile, 0, 0); err != nil {
		t.Fatal(err)
	}

	fromData, err := os.ReadFile(fromFile)
	if err != nil {
		t.Fatal(err)
	}
	toData, err := os.ReadFile(toFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(fromData, toData) {
		t.Error("Files do not match")
	}
}

func TestCopyWithOffsetAndLimit(t *testing.T) {
	fromFile := filepath.Join("testdata", "input.txt")
	toFile := filepath.Join(t.TempDir(), "copy.txt")

	defer func(name string) {
		_ = os.Remove(name)
	}(toFile)

	if err := Copy(fromFile, toFile, 10, 110); err != nil {
		t.Fatal(err)
	}

	fromData, err := os.ReadFile(fromFile)
	if err != nil {
		t.Fatal(err)
	}
	toData, err := os.ReadFile(toFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(fromData[10:120], toData) {
		t.Error("Files do not match")
	}
}

func TestOffsetExceedsFileSize(t *testing.T) {
	fromFile := filepath.Join("testdata", "input.txt")
	toFile := filepath.Join(t.TempDir(), "copy.txt")

	defer func(name string) {
		_ = os.Remove(name)
	}(toFile)

	err := Copy(fromFile, toFile, math.MaxInt64, 110)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !errors.Is(err, ErrOffsetExceedsFileSize) {
		t.Errorf("Expected error %q, got %q", ErrOffsetExceedsFileSize, err.Error())
	}
}

func TestNonRegularFile(t *testing.T) {
	dir := t.TempDir()
	toFile := filepath.Join(t.TempDir(), "copy.txt")

	defer func(name string) {
		_ = os.Remove(name)
	}(toFile)

	err := Copy(dir, toFile, math.MaxInt64, 110)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !errors.Is(err, ErrUnsupportedFile) {
		t.Errorf("Expected error %q, got %q", ErrUnsupportedFile, err.Error())
	}
}

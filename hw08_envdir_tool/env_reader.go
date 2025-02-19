package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("error accessing directory %s: %w", dir, err)
	}

	if !dirInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", dir, err)
	}

	environment := make(Environment)
	for _, dirEntry := range dirEntries {
		if err := process(dir, dirEntry, environment); err != nil {
			return nil, fmt.Errorf("error processing directory entry %s: %w", dirEntry.Name(), err)
		}
	}

	return environment, nil
}

func process(dir string, dirEntry os.DirEntry, environment Environment) error {
	if dirEntry.IsDir() {
		return nil
	}

	if strings.Contains(dirEntry.Name(), "=") {
		return nil
	}

	path := filepath.Join(dir, dirEntry.Name())
	fileInfo, err := dirEntry.Info()
	if err != nil {
		return fmt.Errorf("error reading file info %s: %w", path, err)
	}

	if fileInfo.Size() == 0 {
		environment[fileInfo.Name()] = EnvValue{Value: "", NeedRemove: true}
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", path, err)
	}
	defer safeCloser{file, path}.safeClose()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return fmt.Errorf("error reading file %s: %w", path, scanner.Err())
	}

	value := scanner.Text()
	value = strings.ReplaceAll(value, "\x00", "\n")
	value = strings.TrimRight(value, " \t\n")

	environment[fileInfo.Name()] = EnvValue{Value: value, NeedRemove: false}
	return nil
}

type safeCloser struct {
	io.Closer
	name string
}

func (sf safeCloser) safeClose() {
	err := sf.Close()
	if err != nil {
		log.Printf("error closing %s: %v", sf.name, err)
	}
}

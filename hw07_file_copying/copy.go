package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	inputFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error opening file to read from: %w", err)
	}
	defer SafeCloser{inputFile, "file to read from"}.SafeClose()

	fi, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("error getting file to read from info: %w", err)
	}

	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	size := fi.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	if maxLimit := size - offset; limit == 0 || limit >= maxLimit {
		limit = maxLimit
	}

	outputFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error creating file to write to: %w", err)
	}
	defer SafeCloser{outputFile, "file to write to"}.SafeClose()

	if _, err = inputFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking file to read from: %w", err)
	}

	bar := pb.Start64(limit)
	defer bar.Finish()

	reader := bar.NewProxyReader(io.LimitReader(inputFile, limit))

	if _, err = io.Copy(outputFile, reader); err != nil {
		return fmt.Errorf("error copying data: %w", err)
	}

	return nil
}

type SafeCloser struct {
	io.Closer
	name string
}

func (sf SafeCloser) SafeClose() {
	err := sf.Close()
	if err != nil {
		log.Printf("error closing %s: %v", sf.name, err)
	}
}

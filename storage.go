package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type PathTransformFunc func(string) string

var DefaultPathTransformFunc = func(path string) string {
	return path
}

type StorageOpts struct {
	PathTransformFunc PathTransformFunc
}
type Storage struct {
	StorageOpts
}

func NewStorage(opts StorageOpts) *Storage {
	return &Storage{opts}
}

// writeFileStream writes data from an io.Reader to a file identified by the key.

func (s *Storage) writeFileStream(key string, reader io.Reader) error {
	pathName := s.PathTransformFunc(key)
	filename := "somefilename"
	filePath := filepath.Join(pathName, filename)
	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	result, err := io.Copy(file, reader)
	if err != nil {
		return err
	}
	log.Printf("copied %d bytes to disk: %s \n", result, pathName)
	return nil
}

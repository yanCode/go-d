package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type PathTransformFunc func(string) string

func CasPathTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	block_size := 5
	sliceLen := len(hashStr) / block_size
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*block_size, (i+1)*block_size
		paths[i] = hashStr[from:to]
	}
	return strings.Join(paths, "/")
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

func (s *Storage) writeStream(key string, reader io.Reader) error {
	pathName := s.PathTransformFunc(key)
	filename := "somefilename"
	filePath := filepath.Join(pathName, filename)
	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
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

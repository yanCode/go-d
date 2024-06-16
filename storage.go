package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type PathTransformFunc func(string) Pathkey

type Pathkey struct {
	Pathname string
	Original string
}

func CasPathTransformFunc(key string) Pathkey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	block_size := 5
	sliceLen := len(hashStr) / block_size
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*block_size, (i+1)*block_size
		paths[i] = hashStr[from:to]
	}
	//return strings.Join(paths, "/")
	return Pathkey{
		Original: hashStr,
		Pathname: strings.Join(paths, "/"),
	}
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
	if err := os.MkdirAll(pathName.Pathname, os.ModePerm); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, reader)
	if err != nil {
		return err
	}
	filenameBytes := md5.Sum(buf.Bytes())
	filename := hex.EncodeToString(filenameBytes[:])
	file, err := os.Create(filepath.Join(pathName.Pathname, filename))

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

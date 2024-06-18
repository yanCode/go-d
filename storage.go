package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FirstPathName() string {
	slice := strings.Split(p.PathName, "/")
	if len(slice) == 0 {
		return ""
	}
	return slice[0]
}
func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

func CasPathTransformFunc(key string) PathKey {
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
	return PathKey{
		FileName: hashStr,
		PathName: strings.Join(paths, "/"),
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
func (s *Storage) Read(key string) (io.Reader, error) {
	file, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	return buffer, err

}
func (s *Storage) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	return os.Remove(pathKey.PathName)
}
func (s *Storage) readStream(key string) (*os.File, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(pathKey.PathName)
}

// writeFileStream writes data from an io.Reader to a file identified by the key.

func (s *Storage) writeStream(key string, reader io.Reader) error {
	pathName := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathName.PathName, os.ModePerm); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, reader)
	if err != nil {
		return err
	}
	filenameBytes := md5.Sum(buf.Bytes())
	filename := hex.EncodeToString(filenameBytes[:])
	file, err := os.Create(filepath.Join(pathName.PathName, filename))

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

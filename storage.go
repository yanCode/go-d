package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github/yanCode/go-d/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultRootFolderName = "d-storage"

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

	block_size := 16
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

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type StorageOpts struct {
	RootDir           string
	PathTransformFunc PathTransformFunc
	ListenAddr        string //this is used to debug the server address
}
type Storage struct {
	StorageOpts
}

func NewStorage(opts StorageOpts) *Storage {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.RootDir) == 0 {
		opts.RootDir = defaultRootFolderName
	}
	return &Storage{opts}
}
func (s *Storage) Read(id string, key string) (int64, io.Reader, error) {
	return s.readStream(id, key)
}

func (s *Storage) Delete(id string, key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()
	firstPathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.RootDir, id, pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *Storage) readStream(id string, key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.RootDir, id, pathKey.FullPath())
	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}
	file_info, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}
	return file_info.Size(), file, nil
}

func (s *Storage) Has(id string, key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := filepath.Join(s.RootDir, id, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Storage) openFileForWriting(id string, key string) (*os.File, error) {
	pathkey := s.PathTransformFunc(key)
	pathNameWithRoot := filepath.Join(s.RootDir, id, pathkey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return nil, err
	}
	fullPathWithRoot := filepath.Join(s.RootDir, id, pathkey.FullPath())
	utils.Logger.Printf("Server[%s] is creating a file to write in: %s\n", s.ListenAddr, fullPathWithRoot)
	return os.Create(fullPathWithRoot)
}

func (s *Storage) WriteDecrypt(encKey []byte, id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	n, err := copyDecrypt(encKey, r, f)
	return int64(n), err
}

// writeFileStream writes data from an io.Reader to a file identified by the key.

func (s *Storage) writeStream(id string, key string, reader io.Reader) (int64, error) {
	file, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return io.Copy(file, reader)

}
func (s *Storage) Write(id string, key string, reader io.Reader) (int64, error) {
	return s.writeStream(id, key, reader)
}

func (s *Storage) Clear() error {
	return os.RemoveAll(s.RootDir)
}

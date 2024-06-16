package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	//a transform always listens and accepts
	opts := StorageOpts{
		PathTransformFunc: func(path string) string {
			return path
		},
	}
	s := NewStorage(opts)
	data := bytes.NewReader([]byte("hello"))
	if err := s.writeStream("my_pic", data); err != nil {
		t.Error(err)
	}

}

func TestStorage(t *testing.T) {
	opts := StorageOpts{
		PathTransformFunc: CasPathTransformFunc,
	}
	s := NewStorage(opts)
	data := bytes.NewReader([]byte("hello"))
	if err := s.writeStream("my_pic", data); err != nil {
		t.Error(err)
	}
}

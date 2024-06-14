package main

import (
	"bytes"
	"testing"
)

func TestStorage(t *testing.T) {
	opts := StorageOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStorage(opts)
	data := bytes.NewReader([]byte("hello"))
	if err := s.writeStream("my_pic", data); err != nil {
		t.Error(err)
	}
}

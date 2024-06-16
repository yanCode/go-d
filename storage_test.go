package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	//a transform always listens and accepts
	expected_path := "5a572/70c08/000be/514f9/e954b/3fa37/e4890/38be1"
	key := "my_test_file"
	pathname := CasPathTransformFunc(key)
	if pathname != expected_path {
		t.Errorf("got %s, want %s", pathname, expected_path)
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

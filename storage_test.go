package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	//a transform always listens and accepts
	expected_path := "5a572/70c08/000be/514f9/e954b/3fa37/e4890/38be1"
	expected_original := "5a57270c08000be514f9e954b3fa37e489038be1"
	key := "my_test_file"
	pathKey := CasPathTransformFunc(key)
	fmt.Println(pathKey.Original)
	if pathKey.Pathname != expected_path {
		t.Errorf("Pathname got %s, want %s", pathKey, expected_path)
	}
	if pathKey.Original != expected_original {
		t.Errorf("Original got %s, want %s", pathKey.Original, expected_original)
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

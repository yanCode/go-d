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
	fmt.Println(pathKey.FileName)
	if pathKey.PathName != expected_path {
		t.Errorf("PathName got %s, want %s", pathKey, expected_path)
	}
	if pathKey.FileName != expected_original {
		t.Errorf("FileName got %s, want %s", pathKey.FileName, expected_original)
	}

}

func TestStorage(t *testing.T) {

	opts := StorageOpts{
		PathTransformFunc: CasPathTransformFunc,
		RootDir:           "/Users/y/drills/go-d-system/assets",
	}
	s := NewStorage(opts)
	defer teardown(t, s)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("key_%d", i)
		data := []byte("hello world jpg bytes of number " + fmt.Sprintf("%d", i))
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}
		if ok := s.Has(key); !ok {
			t.Errorf("key %s should exist", key)
		}
		_, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}
	}
}
func TestStorage_Delete(t *testing.T) {
	opts := StorageOpts{
		PathTransformFunc: CasPathTransformFunc,
	}
	s := NewStorage(opts)
	key := "my_pic"
	data := []byte("hello world")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}
func teardown(t *testing.T, s *Storage) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

func TestStorage_Read(t *testing.T) {
	opts := StorageOpts{
		PathTransformFunc: CasPathTransformFunc,
	}
	s := NewStorage(opts)
	key := "my_pic"
	data := []byte("hello world")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if _, err := s.Read(key); err != nil {
		t.Error(err)
	}
}

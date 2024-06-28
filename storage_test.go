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
	s := newStore()
	id := generateId()
	defer tearDown(t, s)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("key_%d", i)
		data := []byte("Some photo bytes")
		if _, err := s.writeStream(id, key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}
		if ok := s.Has(id, key); !ok {
			t.Errorf("key %s not found", key)
		}
		_, reader, err := s.Read(id, key)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(reader)
	}
}
func newStore() *Storage {
	opts := StorageOpts{
		PathTransformFunc: CasPathTransformFunc,
	}
	return NewStorage(opts)
}
func tearDown(t *testing.T, s *Storage) {
	if err := s.Clear(); err != nil {
		t.Errorf("failed to clear storage")
	}
}

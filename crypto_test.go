package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {
	payload := "hello world nice to see you"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	if _, err := copyEncrypt(key, src, dst); err != nil {
		t.Error(err)
	}
	fmt.Println(len(payload))
	fmt.Println(len(dst.String()))

	out := new(bytes.Buffer)
	nw, err := copyDecrypt(key, dst, out)
	if err != nil {
		t.Error(err)
	}
	if nw != 16+len(payload) {
		t.Fail()
	}
	if out.String() != payload {
		t.Errorf("decryption failed!!!")
	}
}

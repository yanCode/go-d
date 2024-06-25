package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {
	payload := "hello world"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	if _, err := copyEncrypt(key, src, dst); err != nil {
		t.Error(err)
	}
	fmt.Println(len(payload))
	fmt.Println(len(dst.String()))

}

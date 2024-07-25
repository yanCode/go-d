package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dst io.Writer) (int, error) {
	var (
		buf      = make([]byte, 32*1024)
		totalLen = blockSize
	)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			nn, err := dst.Write(buf[:n])
			if err != nil {
				return 0, err
			}
			totalLen += nn
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return totalLen, nil
}

func copyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize()) // 16 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// prepend the IV to the file.
	if _, err := dst.Write(iv); err != nil {
		return 0, err
	}
	//if _, err := src.Read(iv); err != nil {
	//	return 0, err
	//}
	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}

func newEncryptionKey() []byte {
	keyBuf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, keyBuf)
	if err != nil {
		return nil
	}
	return keyBuf
}
func generateId() string {
	buffer := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buffer)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(buffer)
}
func hashKey(key string) string {
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func copyDecrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}
	// Read the IV from the given io.Reader which, in our case should be the
	// the block.BlockSize() bytes we read.
	iv := make([]byte, block.BlockSize())
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}
	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}

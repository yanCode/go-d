package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	//a transport always listens and accepts
	listAddr := "127.0.0.1:12345"
	tr := NewTCPTransport(listAddr)
	assert.Equal(t, tr.listenAddress, listAddr)

}

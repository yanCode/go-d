package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	//a transport always listens and accepts
	opts := TCPTransportOptions{
		ListenAddr:    ":3000",
		HandshakeFunc: NopHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddr, ":3000")
	assert.Nil(t, tr.ListenAddAccept())

}

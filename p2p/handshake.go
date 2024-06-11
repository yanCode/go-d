package p2p

import "fmt"

var ErrInvalidHandshake = fmt.Errorf("invalid handshake")

type HandshakeFunc func(Peer) error

func NopHandshakeFunc(peer Peer) error {
	return nil
}

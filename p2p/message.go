package p2p

import "net"

type Rpc struct {
	From    net.Addr
	Payload []byte
}

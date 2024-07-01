package p2p

import "net"

/*
*
Peer is an interface that represents any remote node in the network.
*/
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

/*
*
Transport is an interface that represents the communication between the nodes in the networks.
it can be many forms, like (TCP, UDP, websockets, ...).
*/
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan Rpc
	Close() error
}

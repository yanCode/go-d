package p2p

/*
*
Peer is an interface that represents any remote node in the network.
*/
type Peer interface {
	Close() error
}

/*
*
Transport is an interface that represents the communication between the nodes in the networks.
it can be many forms, like (TCP, UDP, websockets, ...).
*/
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan Rpc
	Close() error
}

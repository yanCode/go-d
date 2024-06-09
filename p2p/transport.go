package p2p

type Peer interface{}
type Transport interface {
	ListenAndAccept() error
	// Listen(Peer) error
	// Accept() (Peer, error)
	// Connect(Peer) error
}

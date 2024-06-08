package p2p

type Peer interface{}
type Transport interface {
	// Listen(Peer) error
	// Accept() (Peer, error)
	// Connect(Peer) error
}

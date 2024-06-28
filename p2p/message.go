package p2p

const (
	IncomingStream = 0x1
	OutgoingStream = 0x2
)

type Rpc struct {
	From    string
	Payload []byte
	Stream  bool
}

package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	mu            sync.RWMutex
	peers         map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {

	return &TCPTransport{
		listenAddress: listenAddr,
	}
}
func (t *TCPTransport) ListenAddAccept() error {
	ln, err := net.Listen("tcp", t.listenAddress)
	t.listener = ln
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	return nil
}
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %v\n", err)
		}
		go t.handleConn(conn)
	}
}
func (t *TCPTransport) handleConn(conn net.Conn) {
	fmt.Printf("new incomming connection: %v\n", conn.RemoteAddr())
}

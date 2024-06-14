package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn net.Conn
	//if we dial and retrieve a connection, then outbound is true
	//if we accept and retrieve a connection, then outbound is false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOptions struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener
	decoder  Decoder
	rpcCh    chan Rpc
	mu       sync.RWMutex
	peers    map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
		rpcCh:               make(chan Rpc),
	}
}
func (t *TCPTransport) Consume() <-chan Rpc {
	return t.rpcCh
}
func (t *TCPTransport) ListenAddAccept() error {
	ln, err := net.Listen("tcp", t.ListenAddr)
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
		fmt.Printf("new incomming connection: %v\n", conn)
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("new peer: %v\n", peer)
	//if err := t.HandshakeFunc(peer); err != nil {
	//	conn.Close()
	//	fmt.Printf("TCP handshake error: %v\n", err)
	//	return
	//}
	rpc := &Rpc{}
	//buf := make([]byte, 2000)
	for {
		//if err := t.Decoder.Decode(conn, rpc); err != nil {
		//	fmt.Printf("TCP error decoding message: %v\n", err)
		//	continue
		//}

		if err := t.Decoder.Decode(conn, rpc); err != nil {
			fmt.Printf("TCP error reading message: %v\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		fmt.Printf("got message: %+v\n", rpc)
	}

}

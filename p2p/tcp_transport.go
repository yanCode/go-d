package p2p

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type TcpPeer struct {
	Conn net.Conn
	//if we dial and retrieve a connection, then outbound is true
	//if we accept and retrieve a connection, then outbound is false
	outbound bool
}

func (t *TcpPeer) Send(b []byte) error {
	_, err := t.Conn.Write(b)
	return err
}

func (t *TcpPeer) Close() error {
	return t.Conn.Close()
}

func NewTcpPeer(conn net.Conn, outbound bool) *TcpPeer {
	return &TcpPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOptions struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
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
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %v\n", err)
		}
		fmt.Printf("new incomming connection: %v\n", conn)
		go t.handleConn(conn, false)
	}
}
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil
	}
	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) handleConn(conn net.Conn, isOutbound bool) {
	//var err error
	defer func() {
		fmt.Printf("closing connection: %v\n", conn)
		conn.Close()
	}()
	peer := NewTcpPeer(conn, isOutbound)
	fmt.Printf("new peer: %v\n", peer)
	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %v\n", err)
		return
	}
	if t.OnPeer != nil {
		//if err = t.OnPeer(peer); err != nil {
		//	conn.Close()
		//	fmt.Printf("OnPeer error: %v\n", err)
		//	return
		//}
	}
	rpc := Rpc{}
	for {

		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP error reading message: %v\n", err)
			return
		}
		rpc.From = conn.RemoteAddr()
		t.rpcCh <- rpc
		fmt.Printf("got message: %+v\n", rpc)
	}

}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

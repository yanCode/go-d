package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type TcpPeer struct {
	net.Conn
	//if we dial and retrieve a connection, then outbound is true
	//if we accept and retrieve a connection, then outbound is false
	outbound  bool
	waitGroup *sync.WaitGroup
}

func (t *TcpPeer) CloseStream() {
	t.waitGroup.Done()
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
		Conn:      conn,
		outbound:  outbound,
		waitGroup: &sync.WaitGroup{},
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
	rpcCh    chan Rpc
}

func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
		rpcCh:               make(chan Rpc, 1024),
	}
}

func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

func (t *TCPTransport) Consume() <-chan Rpc {
	return t.rpcCh
}
func (t *TCPTransport) ListenAddAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	log.Printf("TCP transport listening on port: %s\n", t.ListenAddr)
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
func (t *TCPTransport) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) handleConn(conn net.Conn, isOutbound bool) {
	var err error
	defer func() {
		fmt.Printf("closing connection: %v\n", conn)
		conn.Close()
	}()
	peer := NewTcpPeer(conn, isOutbound)
	fmt.Printf("new peer: %v\n", peer)
	if err = t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %v\n", err)
		return
	}
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			conn.Close()
			fmt.Printf("OnPeer error: %v\n", err)
			return
		}
	}

	for {
		rpc := Rpc{}
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP error reading message: %v\n", err)
			return
		}
		rpc.From = conn.RemoteAddr().String()
		if rpc.Stream {
			peer.waitGroup.Add(1)
			fmt.Printf("[%s] incoming stream, waiting...\n", conn.RemoteAddr())
			peer.waitGroup.Wait()
			fmt.Printf("[%s] stream closed, resuming read loop\n", conn.RemoteAddr())
			continue
		}
		t.rpcCh <- rpc

	}

}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

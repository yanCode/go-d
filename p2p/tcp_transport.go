package p2p

import (
	"errors"
	"fmt"
	"github/yanCode/go-d/utils"
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
	addr := t.ListenAddr
	utils.Logger.Printf("Server:[%s] TCP transport listening on port: %s\n", addr, addr)
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			utils.Logger.Printf("Server[%s]: TCP accept error: %v\n", t.ListenAddr, err)
		}
		utils.Logger.Printf("Server[%s]: accept new incoming connection from: %s\n", t.ListenAddr, conn.RemoteAddr())
		go t.handleConn(conn, false)
	}
}
func (t *TCPTransport) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	utils.Logger.Printf("Server: [%s]  New outgoing connection dialed from: %s\n", t.ListenAddr, conn.RemoteAddr())
	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) handleConn(conn net.Conn, isOutbound bool) {
	var err error
	defer func() {
		utils.Logger.Printf("closing connection: %v\n", conn)
		conn.Close()
	}()
	peer := NewTcpPeer(conn, isOutbound)
	if err = t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %v\n", err)
		return
	}
	if t.OnPeer != nil {
		utils.Logger.Printf("server[ %s]: add a new peer from: %s\n", t.Addr(), peer.RemoteAddr())
		if err = t.OnPeer(peer); err != nil {
			conn.Close()
			utils.Logger.Printf("OnPeer error: %v\n", err)
			return
		}
	}

	for {
		rpc := Rpc{}
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			utils.Logger.Printf("[%s] TCP error reading message:  %v\n", err)
			return
		}
		utils.Logger.Printf("server: [ %s] A peer from: %s accepted a PRC of which steam = %t... \n", t.ListenAddr, peer.RemoteAddr(), rpc.Stream)
		//utils.Logger.Printf("%s", string(rpc.Payload))
		rpc.From = conn.RemoteAddr().String()
		if rpc.Stream {
			peer.waitGroup.Add(1)
			utils.Logger.Printf("[%s] there is an incoming stream, waiting...\n", conn.RemoteAddr())
			peer.waitGroup.Wait()
			utils.Logger.Printf("[%s] the stream closed, resuming read loop\n", conn.RemoteAddr())
			continue
		}
		t.rpcCh <- rpc

	}

}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

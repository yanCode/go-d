package main

import (
	"fmt"
	"github/yanCode/go-d/p2p"
	"log"
)

func OnPeer(peer p2p.Peer) error {
	//peer.Close()
	fmt.Println("got peer, doing some logic with peer outside of TCPTransport: ", peer)
	return nil
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tpcTransportOptions := p2p.TCPTransportOptions{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NopHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer, //todo
	}
	tcpTransport := p2p.NewTCPTransport(tpcTransportOptions)

	fileServerOptions := FileServerOptions{
		StorageRoot:       "/Users/y/drills/go-d-system/" + listenAddr + "_3000",
		PathTransformFunc: CasPathTransformFunc,
		Transport:         *tcpTransport,
		BootstrapNodes:    nodes,
	}
	return NewFileServer(fileServerOptions)
}

func main() {

	s1 := makeServer(":3000", ":3000")
	//go func() {
	log.Fatal(s1.Start())
	//}()
	//s2 := makeServer(":4000", ":3000")
	//log.Fatal(s2.Start())
}

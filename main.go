package main

import (
	"fmt"
	"github/yanCode/go-d/p2p"
	"log"
	"time"
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
	}
	tcpTransport := p2p.NewTCPTransport(tpcTransportOptions)

	fileServerOptions := FileServerOptions{
		EncKey:            newEncryptionKey(),
		StorageRoot:       "/Users/y/drills/go-d-system/" + listenAddr + "_3000",
		PathTransformFunc: CasPathTransformFunc,
		Transport:         *tcpTransport,
		BootstrapNodes:    nodes,
	}
	return NewFileServer(fileServerOptions)
}

func main() {
	s1 := makeServer(":8001", "")
	s2 := makeServer(":8007", "")
	s3 := makeServer("8005", ":8001", ":8007")
	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(500 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()
	time.Sleep(2 * time.Second)
	s3.Start()
	time.Sleep(2 * time.Second)
	select {}
}

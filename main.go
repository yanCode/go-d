package main

import (
	"fmt"
	"github/yanCode/go-d/p2p"
	"log"
	"time"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("got peer, doing some logic with peer outside of TCPTransport: ", peer)
	return nil
}

func main() {
	tpcTransportOptions := p2p.TCPTransportOptions{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NopHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer, //todo
	}
	tcpTransport := p2p.NewTCPTransport(tpcTransportOptions)

	fileServerOptions := FileServerOptions{
		StorageRoot:       "/Users/y/drills/go-d-system/assets_3000",
		PathTransformFunc: CasPathTransformFunc,
		Transport:         *tcpTransport,
	}
	s := NewFileServer(fileServerOptions)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()
}

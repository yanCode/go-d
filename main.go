package main

import (
	"fmt"
	"github/yanCode/go-d/p2p"
	"log"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("got peer, doing some logic with peer outside of TCPTransport: ", peer)
	return nil
}

func main() {
	tcpOptions := p2p.TCPTransportOptions{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NopHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOptions)

	go func() {
		for {
			rpc := <-tr.Consume()
			log.Printf("got message: %+v\n", rpc)
		}
	}()
	if err := tr.ListenAddAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}

package main

import (
	"github/yanCode/go-d/p2p"
	"log"
)

func main() {
	tcpOptions := p2p.TCPTransportOptions{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NopHandshakeFunc,
		Decoder:       p2p.GobDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOptions)
	if err := tr.ListenAddAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}

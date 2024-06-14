package main

import (
	"github/yanCode/go-d/p2p"
	"log"
)

func main() {
	tcpOptions := p2p.TCPTransportOptions{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NopHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
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

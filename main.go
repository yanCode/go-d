package main

import (
	"github/yanCode/go-d/p2p"
	"log"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")
	if err := tr.ListenAddAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}

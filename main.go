package main

import (
	"bytes"
	"fmt"
	"github/yanCode/go-d/p2p"
	"io/ioutil"
	"log"
	"time"
)

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
	s := NewFileServer(fileServerOptions)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {
	s1 := makeServer(":8001", "")
	s2 := makeServer(":8002", "")
	s3 := makeServer(":8003", ":8001", ":8002")
	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(500 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()
	time.Sleep(2 * time.Second)
	go func() {
		log.Fatal(s3.Start())
	}()
	time.Sleep(2 * time.Second)

	for i := 0; i < 1; i++ {
		key := fmt.Sprintf("picture_%d.png", i)
		data := bytes.NewReader([]byte("my big data file here!"))
		s3.Store(key, data)
		if err := s3.storage.Delete(s3.ID, key); err != nil {
			log.Fatal(err)
		}
		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}

package main

import (
	"bytes"
	"fmt"
	"github/yanCode/go-d/p2p"
	"github/yanCode/go-d/utils"
	"io"
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
		StorageRoot:       "/Users/y/drills/go-d-system/" + listenAddr[1:] + "_network",
		PathTransformFunc: CasPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}
	s := NewFileServer(fileServerOptions)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {
	utils.Logger.Println("start to mimic  2 file servers ....")
	s1 := makeServer(":3000", "")
	//s2 := makeServer(":8002", "")
	s3 := makeServer(":4000", ":3000")
	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(500 * time.Millisecond)
	//go func() { log.Fatal(s2.Start()) }()
	time.Sleep(2 * time.Second)
	go func() {
		log.Fatal(s3.Start())
	}()
	time.Sleep(2 * time.Second)
	for i := 0; i < 10; i++ {
		//i := 1
		key := fmt.Sprintf("picture_%d.png", i)
		data := bytes.NewReader([]byte("my big data file here! which is very very big"))
		err := s3.Store(key, data)
		if err != nil {
			log.Fatal(err)
		}

		if err := s3.storage.Delete(s3.ID, key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		result, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(result))
	}
}

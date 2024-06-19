package main

import (
	"fmt"
	"github/yanCode/go-d/p2p"
)

type FileServer struct {
	FileServerOptions
	storage *Storage
	quitch  chan struct{}
}
type FileServerOptions struct {
	ListenAddr          string
	StorageRoot         string
	PathTransformFunc   PathTransformFunc
	Transport           p2p.TCPTransport
	TCPTransportOptions p2p.TCPTransportOptions
}

func NewFileServer(opts FileServerOptions) *FileServer {
	storeOpts := StorageOpts{
		RootDir:           opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOptions: opts,
		storage:           NewStorage(storeOpts),
		quitch:            make(chan struct{}),
	}
}

func (s *FileServer) Start() error {
	fmt.Printf("[%s] starting fileserver...\n", "todo")
	if err := s.Transport.ListenAddAccept(); err != nil {
		return err
	}
	go s.loop()
	return nil
}
func (s *FileServer) loop() {
	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitch:
			return
		}
	}
}
func (s *FileServer) Stop() {
	fmt.Printf("[%s] stopping fileserver...\n", "todo")
	close(s.quitch)
}

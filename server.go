package main

import (
	"fmt"
	"github/yanCode/go-d/p2p"
	"log"
)

type FileServer struct {
	FileServerOptions `json:"file_server_options"`
	storage           *Storage      `json:"storage,omitempty"`
	quitCh            chan struct{} `json:"quitCh,omitempty"`
}
type FileServerOptions struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.TCPTransport
	BootstrapNodes    []string
}

func NewFileServer(opts FileServerOptions) *FileServer {
	storeOpts := StorageOpts{
		RootDir:           opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOptions: opts,
		storage:           NewStorage(storeOpts),
		quitCh:            make(chan struct{}),
	}
}

func (s *FileServer) Start() error {
	fmt.Printf("[%s] starting fileserver...\n", "todo")
	if err := s.Transport.ListenAddAccept(); err != nil {
		return err
	}
	s.loop()
	return nil
}
func (s *FileServer) loop() {
	defer func() {
		s.Transport.Close()
		log.Println("file server stopped for user's quit action..")
	}()
	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitCh:
			return
		}
	}
}
func (s *FileServer) Stop() {
	close(s.quitCh)
	fmt.Printf("[%s] stopping fileserver...\n", "todo")
}

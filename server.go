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
func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		go func(addr string) {
			//fmt.Printf("attemping to connect to connect with remote", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error:", err)
				panic(err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {
	fmt.Printf("[%s] starting fileserver...\n", "todo")
	if err := s.Transport.ListenAddAccept(); err != nil {
		return err
	}
	if len(s.BootstrapNodes) > 0 {
		if err := s.bootstrapNetwork(); err != nil {
			return err
		}
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

func (s *FileServer) Has(key string) bool {
	return s.storage.Has(key)
}

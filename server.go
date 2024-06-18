package main

import "github/yanCode/go-d/p2p"

type FileServer struct {
	FileServerOptions
	storage *Storage
}
type FileServerOptions struct {
	ListenAddr          string
	StorageRoot         string
	PathTransformFunc   PathTransformFunc
	Transport           p2p.Transport
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
	}
}

func (s *FileServerOptions) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	return nil
}

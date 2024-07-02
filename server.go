package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github/yanCode/go-d/p2p"
	"io"
	"log"
	"sync"
	"time"
)

type FileServer struct {
	FileServerOptions
	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	storage  *Storage
	quitCh   chan struct{}
}

type FileServerOptions struct {
	ID                string
	EncKey            []byte
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
	if len(opts.ID) == 0 {
		opts.ID = generateId()
	}
	return &FileServer{
		FileServerOptions: opts,
		storage:           NewStorage(storeOpts),
		quitCh:            make(chan struct{}),
	}
}
func (s *FileServer) broadcast(message *Message) error {
	buffer := new(bytes.Buffer)
	if err := gob.NewEncoder(buffer).Encode(message); err != nil {
		panic(err)
	}
	for _, peer := range s.peers {
		peer.Send([]byte{p2p.IncomingStream})
		if err := peer.Send(buffer.Bytes()); err != nil {
			return err
		}
	}
	return nil

}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		go func(addr string) {
			fmt.Printf("[%s] attempting to connect to  with remote: %s\n", s.Transport.Addr(), addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error:", err)
				panic(err)
			}
			fmt.Println("connected to remote: ", addr)
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {
	fmt.Printf("[%s] starting fileserver...\n", s.Transport.Addr())
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
func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.storage.Has(s.ID, key) {
		fmt.Printf("[%s] serving file (%s) from local disk\n", s.Transport.Addr(), key)
		_, r, err := s.storage.Read(s.ID, key)
		return r, err
	}
	fmt.Printf("[%s] serving file (%s) from network\n", s.Transport.Addr(), key)
	msg := Message{Payload: MessageGetFile{
		ID:  s.ID,
		Key: hashkey(key),
	}}
	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}
	time.Sleep(time.Millisecond * 500)
	for _, peer := range s.peers {
		var fileSize int64
		binary.Read(peer, binary.LittleEndian, &fileSize)
		n, err := s.storage.WriteDecrypt(s.EncKey, s.ID, key, io.LimitReader(peer, fileSize))
		if err != nil {
			return nil, err
		}
		fmt.Printf("[%s] received (%d) bytes over the network from (%s)", s.Transport.Addr(), n, peer.RemoteAddr())

		peer.CloseStream()
	}
	return nil, nil
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[peer.RemoteAddr().String()] = peer
	log.Printf("connected with remote %s", peer.RemoteAddr())
	return nil
}

type Message struct {
	Payload any
}
type MessageStoreFile struct {
	ID   string
	Key  string
	Size int64
}
type MessageGetFile struct {
	ID  string
	Key string
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}

func (s *FileServer) Store(key string, reader io.Reader) error {
	var (
		fileBuffer = new(bytes.Buffer)
		tee        = io.TeeReader(reader, fileBuffer)
	)
	size, err := s.storage.Write(s.ID, key, tee)
	if err != nil {
		return err
	}
	msg := Message{
		Payload: MessageStoreFile{
			ID:   s.ID,
			Key:  hashkey(key),
			Size: size + 16,
		},
	}
	if err := s.broadcast(&msg); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 5)
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	mw.Write([]byte{p2p.IncomingStream})
	n, err := copyEncrypt(s.EncKey, fileBuffer, mw)
	fmt.Printf("[%s] received and written (%d) bytes to disk\n", s.Transport.Addr(), n)
	return nil
}

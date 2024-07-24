package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github/yanCode/go-d/p2p"
	"github/yanCode/go-d/utils"
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
	Transport         *p2p.TCPTransport
	BootstrapNodes    []string
}

func NewFileServer(opts FileServerOptions) *FileServer {
	storeOpts := StorageOpts{
		RootDir:           opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
		ListenAddr:        opts.Transport.Addr(), //this is used to debug the server address
	}
	if len(opts.ID) == 0 {
		opts.ID = generateId()
	}
	return &FileServer{
		FileServerOptions: opts,
		storage:           NewStorage(storeOpts),
		quitCh:            make(chan struct{}),
		peers:             make(map[string]p2p.Peer),
	}
}
func (s *FileServer) broadcast(message *Message) error {
	buffer := new(bytes.Buffer)
	if err := gob.NewEncoder(buffer).Encode(message); err != nil {
		panic(err)
	}
	utils.Logger.Printf("[%s] broadcast to all peers\n", s.Transport.Addr())
	for _, peer := range s.peers {
		if err := peer.Send([]byte{p2p.IncomingMessage}); err != nil {
			return err
		}

		if err := peer.Send(buffer.Bytes()); err != nil {
			return err
		}
	}
	return nil

}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			utils.Logger.Printf("[%s] attempting to connect to  remote: %s\n", s.Transport.Addr(), addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error:", err)
				panic(err)
			}
			utils.Logger.Printf("Server {%s} successfully  connected to remote: %s\n ", s.Transport.Addr(), addr)
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {
	//fmt.Printf("Server [%s] starting fileserver...\n", s.Transport.Addr())
	utils.Logger.Printf("Server [%s] starting fileserver...\n", s.Transport.Addr())
	if err := s.Transport.ListenAddAccept(); err != nil {
		return err
	}
	if len(s.BootstrapNodes) > 0 {
		if err := s.bootstrapNetwork(); err != nil {
			return err
		}
	}
	s.handleChanIncoming()
	return nil
}
func (s *FileServer) handleChanIncoming() {
	defer func() {
		err := s.Transport.Close()
		if err != nil {
			return
		}
		log.Println("file server stopped for user's quit action..")
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var message Message

			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&message); err != nil {
				utils.Logger.Println("decoding error: ", err)
			}
			fmt.Printf("Server [%s] from rpcCh get a message is: %+v \n", s.Transport.Addr(), message)
			if err := s.handleMessage(rpc.From, &message); err != nil {
				utils.Logger.Printf("Server [%s] handle message error: ", s.Transport.Addr(), err)
			}

		case <-s.quitCh:
			return
		}
	}
}
func (s *FileServer) Stop() {
	close(s.quitCh)
	fmt.Printf("[%s] stopping fileserver...\n", s.Transport.Addr())
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
		Key: hashKey(key),
	}}
	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}
	time.Sleep(time.Millisecond * 500)
	for _, peer := range s.peers {
		// First read the file size, so we can limit the amount of bytes that we read
		// from the connection, so it will not keep hanging.
		var fileSize int64
		err := binary.Read(peer, binary.LittleEndian, &fileSize)
		if err != nil {
			return nil, err
		}
		n, err := s.storage.WriteDecrypt(s.EncKey, s.ID, key, io.LimitReader(peer, fileSize))
		if err != nil {
			return nil, err
		}
		utils.Logger.Printf("[%s] received (%d) bytes over the network from (%s), about to close stream", s.Transport.Addr(), n, peer.RemoteAddr())
		peer.CloseStream()
	}

	_, r, err := s.storage.Read(s.ID, key)
	return r, err
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[peer.RemoteAddr().String()] = peer
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
	fmt.Printf("Server [%s]: wrote (%d) bytes to disk\n", s.Transport.ListenAddr, size)
	if err != nil {
		return err
	}
	msg := Message{
		Payload: MessageStoreFile{
			ID:   s.ID,
			Key:  hashKey(key),
			Size: size,
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
	fmt.Printf("[%s] sending (%d) bytes to peers\n", s.Transport.Addr(), size)
	_, err = mw.Write([]byte{p2p.IncomingStream})
	if err != nil {
		return err
	}
	//n, err := copyEncrypt(s.EncKey, fileBuffer, mw)
	n, err := io.Copy(mw, fileBuffer)
	utils.Logger.Printf("[%s] Store file Successfully completed, received and written (%d) bytes to disk\n", s.Transport.Addr(), n)
	return nil
}
func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, v)
	case MessageGetFile:
		return s.handleMessageGetFile(from, v)
	}
	return nil
}
func (s *FileServer) handleMessageGetFile(from string, message MessageGetFile) error {
	if !s.storage.Has(message.ID, message.Key) {
		return fmt.Errorf("[%s] need to serve file (%s) but it does not exist on disk", s.Transport.Addr(), message.Key)
	}
	fmt.Printf("[%s] serving file (%s) over the network\n", s.Transport.Addr(), message.Key)
	fileSize, r, err := s.storage.Read(message.ID, message.Key)
	if err != nil {
		return err
	}
	if rc, ok := r.(io.ReadCloser); ok {
		fmt.Println("closing readCloser")
		defer rc.Close()
	}
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) could not be found in the peer list", from)
	}
	peer.Send([]byte{p2p.IncomingStream})
	binary.Write(peer, binary.LittleEndian, fileSize)
	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] written (%d) bytes over the network to %s\n", s.Transport.Addr(), n, from)

	return nil

}

func (s *FileServer) handleMessageStoreFile(from string, message MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) could not be found in the peer list", from)
	}
	n, err := s.storage.Write(message.ID, message.Key, io.LimitReader(peer, message.Size))
	if err != nil {
		return err
	}
	utils.Logger.Printf("[%s] in Storing file, received (%d) bytes over the network from (%s)", s.Transport.Addr(), n, peer.RemoteAddr())
	peer.CloseStream()
	utils.Logger.Printf("[%s] closed stream with (%s)", s.Transport.Addr(), peer.RemoteAddr())
	return nil
}

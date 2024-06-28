package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Rpc) error
}

type GobDecoder struct{}

type DefaultDecoder struct{}

func (dec GobDecoder) Decode(reader io.Reader, msg *Rpc) error {
	return gob.NewDecoder(reader).Decode(msg)

}

func (noop DefaultDecoder) Decode(reader io.Reader, msg *Rpc) error {
	peekBuffer := make([]byte, 1)
	if _, err := reader.Read(peekBuffer); err != nil {
		return nil
	}
	stream := peekBuffer[0] == IncomingStream
	if stream {
		msg.Stream = true
		return nil
	}

	buffer := make([]byte, 1028)
	n, err := reader.Read(buffer)
	if err != nil {
		return err
	}
	msg.Payload = buffer[:n]
	return nil
}

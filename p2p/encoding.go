package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GobDecoder struct{}

type DefaultDecoder struct{}

func (dec GobDecoder) Decode(reader io.Reader, msg *Message) error {
	return gob.NewDecoder(reader).Decode(msg)

}

func (noop DefaultDecoder) Decode(reader io.Reader, msg *Message) error {
	//buf := new(bytes.Buffer)
	buf := make([]byte, 1028)
	n, err := reader.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]
	return nil
}

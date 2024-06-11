package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, any) error
}

type GobDecoder struct{}

func (dec *GobDecoder) Decode(reader io.Reader, msg any) error {
	return gob.NewDecoder(reader).Decode(msg)

}

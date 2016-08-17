package proxy

import (
	"github.com/fagongzi/goetty"
)

var (
	// DECODER TransparentDecoder
	DECODER = NewTransparentDecoder()
	// ENCODER TransparentEncoder
	ENCODER = NewTransparentEncoder()
)

// TransparentDecoder TransparentDecoder
type TransparentDecoder struct {
}

// NewTransparentDecoder create TransparentDecoder
func NewTransparentDecoder() goetty.Decoder {
	return &TransparentDecoder{}
}

// Decode decode
func (d TransparentDecoder) Decode(in *goetty.ByteBuf) (bool, interface{}, error) {
	_, data, err := in.ReadAll()
	return true, data, err
}

// TransparentEncoder TransparentEncoder
type TransparentEncoder struct {
}

// NewTransparentEncoder create TransparentEncoder
func NewTransparentEncoder() goetty.Encoder {
	return &TransparentEncoder{}
}

// Encode Encode
func (e TransparentEncoder) Encode(data interface{}, out *goetty.ByteBuf) error {
	bytes, _ := data.([]byte)
	out.Write(bytes)
	return nil
}

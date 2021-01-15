package bencode

import (
	"io"
	"reflect"
)

// Encoder bencode encoder
type Encoder struct {
	w io.Writer
}

// NewEncoder create encoder
func NewEncoder(w io.Writer) Encoder {
	return Encoder{w}
}

// Encode encode data
func (enc Encoder) Encode(data interface{}) error {
	return encode(enc.w, reflect.ValueOf(data))
}

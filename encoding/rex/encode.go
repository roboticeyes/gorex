package rex

import (
	"io"
)

// Encoder is used to dump a valid REX file buffer into a writer
type Encoder struct {
	w io.Writer
}

// NewEncoder creates a new REX encoder
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode encodes a given REX file buffer into the writer stream.
// The function returns the number of bytes being written to the writer
// and nil if no error occurs.
func (enc *Encoder) Encode(r File) (int, error) {

	return r.Header.Write(enc.w)
}

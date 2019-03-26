package rex

import (
	"fmt"
	"io"
)

// Decoder which can be used to read and decode REX files from a stream
type Decoder struct {
	r   io.Reader
	buf []byte
}

// NewDecoder creates a new REX decoder with a given input stream
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the input from the reader and returns
// a valid REX datastructure.
func (dec *Decoder) Decode(r *File) error {
	return fmt.Errorf("not implemented")
}

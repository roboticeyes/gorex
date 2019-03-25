package rex

import (
	"fmt"
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

	var total int

	n, err := r.Header().Write(enc.w)
	total += n

	// Write PointLists
	for _, p := range r.PointLists {
		n, err = p.Write(enc.w)
		total += p.GetSize() // - totalHeaderSize
		if err != nil {
			return total, err
		}
	}

	// Write Meshes
	for _, m := range r.Meshes {
		n, err = m.Write(enc.w)
		total += m.GetSize() // - totalHeaderSize
		if err != nil {
			return total, err
		}
		fmt.Println("Mesh written ", total)
	}

	return total, nil
}

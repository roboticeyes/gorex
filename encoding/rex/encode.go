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
func (enc *Encoder) Encode(r File) error {

	err := r.Header().Write(enc.w)

	// Write PointLists
	for _, p := range r.PointLists {
		err = p.Write(enc.w)
		if err != nil {
			return err
		}
	}

	// Write Meshes
	for _, m := range r.Meshes {
		err = m.Write(enc.w)
		if err != nil {
			return err
		}
	}

	// Write Materials
	for _, m := range r.Materials {
		err = m.Write(enc.w)
		if err != nil {
			return err
		}
	}

	return nil
}

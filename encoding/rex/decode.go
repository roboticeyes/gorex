package rex

import (
	"encoding/binary"
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
func (dec *Decoder) Decode() (*Header, *File, error) {

	header, err := ReadHeader(dec.r)
	if err != nil {
		return &Header{}, nil, err
	}
	file := &File{}

	for {
		var blockType, version uint16
		var sz uint32
		var id uint64
		if err := binary.Read(dec.r, binary.LittleEndian, &blockType); err != nil {
			if err == io.EOF {
				return header, file, nil
			}
			return header, nil, err
		}
		if err := binary.Read(dec.r, binary.LittleEndian, &version); err != nil {
			return header, nil, err
		}
		if err := binary.Read(dec.r, binary.LittleEndian, &sz); err != nil {
			return header, nil, err
		}
		if err := binary.Read(dec.r, binary.LittleEndian, &id); err != nil {
			return header, nil, err
		}

		// read block
		fmt.Printf("reading type %d version %d sz %d id %d\n", blockType, version, sz, id)
		buf, err := readBlock(dec.r, sz)

		switch blockType {
		case typeMesh:
			mesh, err := ReadMesh(buf)
			mesh.ID = id
			if err != nil {
				file.Meshes = append(file.Meshes, *mesh)
			}
		}

		if err == io.EOF {
			return header, file, nil
		} else if err != nil {
			return header, nil, err
		}
	}
}

func readBlock(r io.Reader, n uint32) ([]byte, error) {
	buf := make([]byte, n)
	if err := binary.Read(r, binary.LittleEndian, &buf); err != nil {
		return nil, err
	}
	return buf, nil
}

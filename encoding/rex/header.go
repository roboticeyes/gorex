package rex

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	// TotalHeaderSize is the number of bytes for each block header
	TotalHeaderSize = 16
)

// Header defines the structure of the REX header
type Header struct {
	magic     [4]byte
	version   uint16
	crc       uint32
	NrBlocks  uint16
	startAddr uint16
	SizeBytes uint64
	reserved  [42]byte
}

// CreateHeader returns a valid fresh header block
func CreateHeader() *Header {
	header := &Header{
		version:   1,
		crc:       0,
		NrBlocks:  0,
		startAddr: 86, // fixed CSB of 22 bytes
		SizeBytes: 0,
	}
	header.magic[0] = 'R'
	header.magic[1] = 'E'
	header.magic[2] = 'X'
	header.magic[3] = '1'
	return header
}

// Write converts the REX header and a dummy CSR and writes it to the given writer
func (h *Header) Write(w io.Writer) (int, error) {
	buf := new(bytes.Buffer)

	var header = []interface{}{
		h.magic,
		h.version,
		h.crc,
		h.NrBlocks,
		h.startAddr,
		h.SizeBytes,
		h.reserved,
		// default CSB block
		uint32(3876),
		uint16(4),
		[]byte("EPSG"),
		float32(0.0),
		float32(0.0),
		float32(0.0),
	}
	for _, v := range header {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return 0, err
		}
	}
	return w.Write(buf.Bytes())
}

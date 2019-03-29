package rex

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	rexFileHeaderSize = 64
	totalHeaderSize   = 16

	typeLineSet          = 0
	typeText             = 1
	typePointList        = 2
	typeMesh             = 3
	typeImage            = 4
	typeMaterial         = 5
	typePeopleSimulation = 6
	typeUnityPackage     = 7
)

// Header defines the structure of the REX header
type Header struct {
	Magic     [4]byte
	Version   uint16
	Crc       uint32
	NrBlocks  uint16
	StartAddr uint16
	SizeBytes uint64
	Reserved  [42]byte
}

// DataBlockHeader stores the header information of a data block
type DataBlockHeader struct {
	Type    uint16
	Version uint16
	Size    uint32
	ID      uint64
}

// CreateHeader returns a valid fresh header block
func CreateHeader() *Header {
	header := &Header{
		Version:   1,
		Crc:       0,
		NrBlocks:  0,
		StartAddr: 86, // fixed CSB of 22 bytes
		SizeBytes: 0,
	}
	header.Magic[0] = 'R'
	header.Magic[1] = 'E'
	header.Magic[2] = 'X'
	header.Magic[3] = '1'
	return header
}

// Write converts the REX header and a dummy CSR and writes it to the given writer
func (h *Header) Write(w io.Writer) (int, error) {
	buf := new(bytes.Buffer)

	var header = []interface{}{
		h.Magic,
		h.Version,
		h.Crc,
		h.NrBlocks,
		h.StartAddr,
		h.SizeBytes,
		h.Reserved,
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

// ReadHeader reads the REX header from a given file
func ReadHeader(r io.Reader) (*Header, error) {

	var header Header
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		fmt.Println("ReadHeader failed:", err)
		return &Header{}, fmt.Errorf("Error during reading header %v", err)
	}

	// read coordinate system block
	var srid uint32
	var sz uint16
	binary.Read(r, binary.LittleEndian, &srid)
	binary.Read(r, binary.LittleEndian, &sz)
	name := make([]byte, sz)
	binary.Read(r, binary.LittleEndian, &name)
	var x, y, z float32
	binary.Read(r, binary.LittleEndian, &x)
	binary.Read(r, binary.LittleEndian, &y)
	binary.Read(r, binary.LittleEndian, &z)

	return &header, nil
}

// ReadDataBlockHeader reads a data block header from reader
func ReadDataBlockHeader(r io.Reader) (DataBlockHeader, error) {
	var hdr DataBlockHeader
	if err := binary.Read(r, binary.LittleEndian, &hdr); err != nil {
		return hdr, err
	}
	return hdr, nil
}

// GetDataBlockHeader returns a new data block header,
// where `sz` denotes the total size of the data block including
// the data block header size (TotalHeaderSize)
func GetDataBlockHeader(blockType, version uint16, blockID uint64, sz int) []byte {

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint16(blockType),
		uint16(version),
		uint32(sz - totalHeaderSize),
		uint64(blockID),
	}

	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			panic(err)
		}
	}

	return buf.Bytes()
}

// String nicely print header
func (h Header) String() string {

	s := fmt.Sprintf("\n")
	s += fmt.Sprintf("| MAGIC          | %-41s |\n", h.Magic)
	s += fmt.Sprintf("| Version        | %-41d |\n", h.Version)
	s += fmt.Sprintf("| CRC            | %-41d |\n", h.Crc)
	s += fmt.Sprintf("| NrBlocks       | %-41d |\n", h.NrBlocks)
	s += fmt.Sprintf("| StartAddr      | %-41d |\n", h.StartAddr)
	s += fmt.Sprintf("| SizeBytes      | %-41d |\n", h.SizeBytes)
	s += fmt.Sprintf("\n")

	return s
}

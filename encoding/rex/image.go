package rex

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	imageBlockVersion = 1
)

const (
	raw24 = iota
	jpeg
	png
)

// Image datastructure
type Image struct {
	ID          uint64
	Compression uint32
	Data        []byte
}

// GetSize returns the estimated size of the block in bytes
func (block *Image) GetSize() int {
	return totalHeaderSize + 4 + len(block.Data)
}

// ReadImage reads a REX image w/o block header
func ReadImage(r io.Reader, hdr DataBlockHeader) (*Image, error) {

	if hdr.Version != imageBlockVersion {
		return nil, fmt.Errorf("Image block version %d is not supported", hdr.Version)
	}
	if hdr.Type != typeImage {
		return nil, fmt.Errorf("Wrong data block type for Image: %d", hdr.Type)
	}

	image := Image{ID: hdr.ID}
	if err := binary.Read(r, binary.LittleEndian, &image.Compression); err != nil {
		fmt.Println("Reading compression failed: ", err)
	}

	image.Data = make([]byte, hdr.Size-4)

	if err := binary.Read(r, binary.LittleEndian, &image.Data); err != nil {
		fmt.Println("Reading Image failed: ", err)
	}

	return &image, nil
}

// Write writes the image including the data header to the given writer
func (block *Image) Write(w io.Writer) error {

	err := WriteDataBlockHeader(w, DataBlockHeader{
		Type:    typeImage,
		Version: imageBlockVersion,
		Size:    uint32(block.GetSize() - totalHeaderSize),
		ID:      block.ID,
	})
	if err != nil {
		return err
	}

	var data = []interface{}{
		block.Compression,
		block.Data,
	}
	for _, v := range data {
		err := binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}
	return nil
}

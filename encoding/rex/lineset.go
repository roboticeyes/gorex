package rex

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	linesetBlockVersion = 1
)

// LineSet stores a list of 3D points which form a polyline. Number of lines are size(p)-1.
// The color is (RGBA)
type LineSet struct {
	ID     uint64
	Colors mgl32.Vec4
	Points []mgl32.Vec3
}

// GetSize returns the estimated size of the block in bytes
func (block *LineSet) GetSize() int {
	return rexDataBlockHeaderSize + 20 + len(block.Points)*12
}

// ReadLineSet reads the block
func ReadLineSet(r io.Reader, hdr DataBlockHeader) (*LineSet, error) {

	var nrVertices uint32
	var ls LineSet

	if err := binary.Read(r, binary.LittleEndian, &ls.Colors); err != nil {
		return nil, fmt.Errorf("Reading failed: %v ", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &nrVertices); err != nil {
		return nil, fmt.Errorf("Reading failed: %v ", err)
	}

	ls.Points = make([]mgl32.Vec3, nrVertices)
	if err := binary.Read(r, binary.LittleEndian, &ls.Points); err != nil {
		return nil, fmt.Errorf("Reading coords failed: %v", err)
	}

	return &ls, nil
}

// Write writes the lineset to the given writer
func (block *LineSet) Write(w io.Writer) error {

	// return if nothing needs to be written
	if len(block.Points) == 0 {
		return nil
	}

	err := WriteDataBlockHeader(w, DataBlockHeader{
		Type:    typeLineSet,
		Version: linesetBlockVersion,
		Size:    uint32(block.GetSize() - rexDataBlockHeaderSize),
		ID:      block.ID,
	})
	if err != nil {
		return err
	}

	var data = []interface{}{
		float32(block.Colors.X()),
		float32(block.Colors.Y()),
		float32(block.Colors.Z()),
		float32(block.Colors.W()),
		uint32(len(block.Points)),
	}
	for _, v := range data {
		err := binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}
	// Points
	for _, p := range block.Points {
		err := binary.Write(w, binary.LittleEndian, p.X())
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, p.Y())
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, p.Z())
		if err != nil {
			return err
		}
	}
	return nil
}

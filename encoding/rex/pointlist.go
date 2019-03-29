package rex

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	pointListBlockVersion = 1
)

// PointList stores a list of (colored) 3D points
type PointList struct {
	ID     uint64
	Points []mgl32.Vec3
	Colors []mgl32.Vec3
}

// GetSize returns the estimated size of the block in bytes
func (block *PointList) GetSize() int {
	return totalHeaderSize + 4 + 4 + len(block.Points)*12 + len(block.Colors)*12
}

// ReadPointList reads a REX material
func ReadPointList(r io.Reader, hdr DataBlockHeader) (*PointList, error) {

	var nrVertices, nrColors uint32

	if err := binary.Read(r, binary.LittleEndian, &nrVertices); err != nil {
		return nil, fmt.Errorf("Reading failed: %v ", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &nrColors); err != nil {
		return nil, fmt.Errorf("Reading failed: %v ", err)
	}

	var pointList PointList

	pointList.Points = make([]mgl32.Vec3, nrVertices)
	if err := binary.Read(r, binary.LittleEndian, &pointList.Points); err != nil {
		return nil, fmt.Errorf("Reading coords failed: %v", err)
	}

	pointList.Colors = make([]mgl32.Vec3, nrColors)
	if err := binary.Read(r, binary.LittleEndian, &pointList.Colors); err != nil {
		return nil, fmt.Errorf("Reading colors failed: %v", err)
	}
	return &pointList, nil
}

// Write writes the pointlist to the given writer
func (block *PointList) Write(w io.Writer) error {

	// return if nothing needs to be written
	if len(block.Points) == 0 {
		return nil
	}

	err := WriteDataBlockHeader(w, DataBlockHeader{
		Type:    typePointList,
		Version: pointListBlockVersion,
		Size:    uint32(block.GetSize() - totalHeaderSize),
		ID:      block.ID,
	})
	if err != nil {
		return err
	}

	var data = []interface{}{
		uint32(len(block.Points)),
		uint32(len(block.Colors)),
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
	// Colors
	for _, c := range block.Colors {
		err := binary.Write(w, binary.LittleEndian, c.X() /* red */)
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, c.Y() /* green */)
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, c.Z() /* blue */)
		if err != nil {
			return err
		}
	}
	return nil
}

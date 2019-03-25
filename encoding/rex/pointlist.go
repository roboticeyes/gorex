package rex

import (
	"bytes"
	"encoding/binary"
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

// Write writes the pointlist to the given writer
func (block *PointList) Write(w io.Writer) (int, error) {

	// return if nothing needs to be written
	if len(block.Points) == 0 {
		return 0, nil
	}

	buf := new(bytes.Buffer)
	var data = []interface{}{
		GetDataBlockHeader(typePointList, pointListBlockVersion, block.ID, block.GetSize()),
		uint32(len(block.Points)),
		uint32(len(block.Colors)),
	}
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return 0, err
		}
	}
	// Points
	for _, p := range block.Points {
		err := binary.Write(buf, binary.LittleEndian, p.X())
		if err != nil {
			return 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, p.Y())
		if err != nil {
			return 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, p.Z())
		if err != nil {
			return 0, err
		}
	}
	// Colors
	for _, c := range block.Colors {
		err := binary.Write(buf, binary.LittleEndian, c.X() /* red */)
		if err != nil {
			return 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, c.Y() /* green */)
		if err != nil {
			return 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, c.Z() /* blue */)
		if err != nil {
			return 0, err
		}
	}
	return w.Write(buf.Bytes())
}

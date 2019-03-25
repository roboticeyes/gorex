package rex

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Point ...
type Point struct {
	X float32
	Y float32
	Z float32
}

// Color ...
type Color struct {
	R float32
	G float32
	B float32
}

// PointList ...
type PointList struct {
	Points []Point
	Colors []Color
}

func getHeader(id uint64, sz int) []byte {

	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint16(2),
		uint16(1),
		uint32(sz - TotalHeaderSize),
		uint64(id),
	}

	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}

// GetSize returns the estimated size of the block in bytes
func (p *PointList) GetSize() int {
	return TotalHeaderSize + 4 + 4 + len(p.Points)*12 + len(p.Colors)*12
}

// Write writes the pointlist to the given writer
func (p *PointList) Write(id uint64, w io.Writer) (int, error) {

	// return if nothing needs to be written
	if len(p.Points) == 0 {
		return 0, nil
	}

	buf := new(bytes.Buffer)

	var data = []interface{}{
		getHeader(id, p.GetSize()),
		uint32(len(p.Points)),
		uint32(len(p.Colors)),
	}
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return 0, err
		}
	}
	// Points
	for _, p := range p.Points {
		err := binary.Write(buf, binary.LittleEndian, p.X)
		if err != nil {
			return 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, p.Y)
		if err != nil {
			return 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, p.Z)
		if err != nil {
			return 0, err
		}
	}
	// Colors
	for _, c := range p.Colors {
		err := binary.Write(buf, binary.LittleEndian, c.R)
		if err != nil {
			return 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, c.G)
		if err != nil {
			return 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, c.B)
		if err != nil {
			return 0, err
		}
	}
	return w.Write(buf.Bytes())
}

package rex

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
		uint32(sz),
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

// Marshal converts a REX point list to a proper buffer
func (p *PointList) Marshal(id uint64) ([]byte, uint64, error) {

	buf := new(bytes.Buffer)
	if len(p.Points) == 0 {
		return buf.Bytes(), 0, fmt.Errorf("No points found")
	}

	sz := 4 + 4 + len(p.Points)*12 + len(p.Colors)*12

	var data = []interface{}{
		getHeader(id, sz),
		uint32(len(p.Points)),
		uint32(len(p.Colors)),
	}
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return buf.Bytes(), 0, err
		}
	}
	// Points
	for _, p := range p.Points {
		err := binary.Write(buf, binary.LittleEndian, p.X)
		if err != nil {
			return buf.Bytes(), 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, p.Y)
		if err != nil {
			return buf.Bytes(), 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, p.Z)
		if err != nil {
			return buf.Bytes(), 0, err
		}
	}
	// Colors
	for _, c := range p.Colors {
		err := binary.Write(buf, binary.LittleEndian, c.R)
		if err != nil {
			return buf.Bytes(), 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, c.G)
		if err != nil {
			return buf.Bytes(), 0, err
		}
		err = binary.Write(buf, binary.LittleEndian, c.B)
		if err != nil {
			return buf.Bytes(), 0, err
		}
	}
	return buf.Bytes(), uint64(sz + 16), nil
}

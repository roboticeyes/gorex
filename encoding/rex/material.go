package rex

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	materialStandardSize = 68
	materialBlockVersion = 1
)

// Material datastructure
type Material struct {
	ID          uint64
	KaRgb       mgl32.Vec3
	KaTextureID uint64
	KdRgb       mgl32.Vec3
	KdTextureID uint64
	KsRgb       mgl32.Vec3
	KsTextureID uint64
	Ns          float32
	Alpha       float32 // 1 is full opaque
}

// NewMaterial creates a new default material
func NewMaterial(id uint64) Material {
	return Material{
		ID:    id,
		KaRgb: mgl32.Vec3{0.8, 0.8, 0.8},
		KdRgb: mgl32.Vec3{0.8, 0.8, 0.8},
		KsRgb: mgl32.Vec3{0.8, 0.8, 0.8},
		Alpha: 1,
	}
}

// GetSize returns the estimated size of the block in bytes
func (block *Material) GetSize() int {
	return rexDataBlockHeaderSize + materialStandardSize
}

// ReadMaterial reads a REX material
func ReadMaterial(r io.Reader, hdr DataBlockHeader) (*Material, error) {

	var rexMaterial struct {
		KaRed, KaGreen, KaBlue float32
		KaTextureID            uint64
		KdRed, KdGreen, KdBlue float32
		KdTextureID            uint64
		KsRed, KsGreen, KsBlue float32
		KsTextureID            uint64
		Ns                     float32
		Alpha                  float32
	}
	if err := binary.Read(r, binary.LittleEndian, &rexMaterial); err != nil {
		return nil, fmt.Errorf("Reading material failed")
	}

	return &Material{
		KaRgb:       mgl32.Vec3{rexMaterial.KaRed, rexMaterial.KaGreen, rexMaterial.KaBlue},
		KaTextureID: rexMaterial.KaTextureID,
		KdRgb:       mgl32.Vec3{rexMaterial.KdRed, rexMaterial.KdGreen, rexMaterial.KdBlue},
		KdTextureID: rexMaterial.KdTextureID,
		KsRgb:       mgl32.Vec3{rexMaterial.KsRed, rexMaterial.KsGreen, rexMaterial.KsBlue},
		KsTextureID: rexMaterial.KsTextureID,
		Ns:          rexMaterial.Ns,
		Alpha:       rexMaterial.Alpha,
	}, nil
}

// Write writes the material to the given writer
func (block *Material) Write(w io.Writer) error {

	err := WriteDataBlockHeader(w, DataBlockHeader{
		Type:    typeMaterial,
		Version: materialBlockVersion,
		Size:    uint32(block.GetSize() - rexDataBlockHeaderSize),
		ID:      block.ID,
	})
	if err != nil {
		return err
	}

	var data = []interface{}{
		float32(block.KaRgb.X()),
		float32(block.KaRgb.Y()),
		float32(block.KaRgb.Z()),
		uint64(block.KaTextureID),

		float32(block.KdRgb.X()),
		float32(block.KdRgb.Y()),
		float32(block.KdRgb.Z()),
		uint64(block.KdTextureID),

		float32(block.KsRgb.X()),
		float32(block.KsRgb.Y()),
		float32(block.KsRgb.Z()),
		uint64(block.KsTextureID),

		float32(block.Ns),
		float32(block.Alpha),
	}
	for _, v := range data {
		err := binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}
	return nil
}

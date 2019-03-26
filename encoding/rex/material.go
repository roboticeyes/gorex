package rex

import (
	"bytes"
	"encoding/binary"
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

// GetSize returns the estimated size of the block in bytes
func (block *Material) GetSize() int {
	return totalHeaderSize + materialStandardSize
}

// Write writes the material to the given writer
func (block *Material) Write(w io.Writer) (int, error) {

	buf := new(bytes.Buffer)

	var data = []interface{}{
		GetDataBlockHeader(typeMaterial, materialBlockVersion, block.ID, block.GetSize()),
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
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return 0, err
		}
	}
	return w.Write(buf.Bytes())
}

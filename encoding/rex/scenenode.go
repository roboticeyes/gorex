package rex

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	sceneNodeSize         = 80
	sceneNodeBlockVersion = 1
)

// SceneNode datastructure
type SceneNode struct {
	ID          uint64
	GeometryID  uint64
	Name        [32]byte
	Translation mgl32.Vec3
	Rotation    mgl32.Vec4
	Scale       mgl32.Vec3
}

// NewSceneNode creates a new empty SceneNode pointing to no geometry
func NewSceneNode(id uint64) SceneNode {
	return SceneNode{
		ID:          id,
		GeometryID:  0,
		Name:        [32]byte{},
		Translation: mgl32.Vec3{0.0, 0.0, 0.0},
		Rotation:    mgl32.Vec4{0.0, 0.0, 0.0, 1.0},
		Scale:       mgl32.Vec3{1.0, 1.0, 1.0},
	}
}

// GetSize returns the estimated size of the block in bytes
func (block *SceneNode) GetSize() int {
	return sceneNodeSize
}

// ReadSceneNode reads the block
func ReadSceneNode(r io.Reader, hdr DataBlockHeader) (*SceneNode, error) {

	var block struct {
		GeometryID     uint64
		Name           [32]byte
		Tx, Ty, Tz     float32
		Rx, Ry, Rz, Rw float32
		Sx, Sy, Sz     float32
	}
	if err := binary.Read(r, binary.LittleEndian, &block); err != nil {
		return nil, fmt.Errorf("Reading SceneNode failed")
	}

	return &SceneNode{
		ID:          hdr.ID,
		GeometryID:  block.GeometryID,
		Name:        block.Name,
		Translation: mgl32.Vec3{block.Tx, block.Ty, block.Tz},
		Rotation:    mgl32.Vec4{block.Rx, block.Ry, block.Rz, block.Rw},
		Scale:       mgl32.Vec3{block.Sx, block.Sy, block.Sz},
	}, nil
}

// Write writes the scenenode
func (block *SceneNode) Write(w io.Writer) error {

	err := WriteDataBlockHeader(w, DataBlockHeader{
		Type:    typeSceneNode,
		Version: sceneNodeBlockVersion,
		Size:    uint32(block.GetSize() - rexDataBlockHeaderSize),
		ID:      block.ID,
	})
	if err != nil {
		return err
	}

	var data = []interface{}{
		uint64(block.GeometryID),
		block.Name,

		float32(block.Translation.X()),
		float32(block.Translation.Y()),
		float32(block.Translation.Z()),

		float32(block.Rotation.X()),
		float32(block.Rotation.Y()),
		float32(block.Rotation.Z()),
		float32(block.Rotation.W()),

		float32(block.Scale.X()),
		float32(block.Scale.Y()),
		float32(block.Scale.Z()),
	}
	for _, v := range data {
		err := binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}
	return nil
}

package rex

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	meshHeaderSize   = 128
	meshBlockVersion = 1
	meshNameMaxSize  = 74
)

// Triangle defines three indices
type Triangle struct {
	V0 uint32
	V1 uint32
	V2 uint32
}

// Mesh datastructure
type Mesh struct {
	ID         uint64
	Name       string
	Coords     []mgl32.Vec3
	Normals    []mgl32.Vec3
	TexCoords  []mgl32.Vec2
	Colors     []mgl32.Vec3
	Triangles  []Triangle
	MaterialID uint64
}

// GetSize returns the estimated size of the block in bytes
func (block *Mesh) GetSize() int {
	return totalHeaderSize + meshHeaderSize +
		len(block.Coords)*12 +
		len(block.Normals)*12 +
		len(block.TexCoords)*8 +
		len(block.Colors)*12 +
		len(block.Triangles)*12
}

// ReadMesh reads a REX mesh
func ReadMesh(buf []byte) (*Mesh, error) {

	r := bytes.NewReader(buf)

	var rexMesh struct {
		Lod, MaxLod                                           uint16
		NrCoords, NrNormals, NrTexCoords, NrColors            uint32
		NrTriangles                                           uint32
		StartCoords, StartNormals, StartTexCoords, StartColor uint32
		StartTriangles                                        uint32
		MaterialID                                            uint64
		NameSize                                              uint16
		Name                                                  [74]byte
	}
	if err := binary.Read(r, binary.LittleEndian, &rexMesh); err != nil {
		fmt.Println("Reading MeshHeader failed: ", err)
	}

	var mesh Mesh
	mesh.Name = string(rexMesh.Name[:rexMesh.NameSize])

	// Read coordinates
	mesh.Coords = make([]mgl32.Vec3, rexMesh.NrCoords)
	if err := binary.Read(r, binary.LittleEndian, &mesh.Coords); err != nil {
		fmt.Println("Reading coords failed: ", err)
	}

	// Read normals
	mesh.Normals = make([]mgl32.Vec3, rexMesh.NrNormals)
	if err := binary.Read(r, binary.LittleEndian, &mesh.Normals); err != nil {
		fmt.Println("Reading normals failed: ", err)
	}

	// Read texture
	mesh.TexCoords = make([]mgl32.Vec2, rexMesh.NrTexCoords)
	if err := binary.Read(r, binary.LittleEndian, &mesh.TexCoords); err != nil {
		fmt.Println("Reading texture failed: ", err)
	}

	// Read color
	mesh.Colors = make([]mgl32.Vec3, rexMesh.NrColors)
	if err := binary.Read(r, binary.LittleEndian, &mesh.Colors); err != nil {
		fmt.Println("Reading colors failed: ", err)
	}

	// Read triangles
	mesh.Triangles = make([]Triangle, rexMesh.NrTriangles)
	if err := binary.Read(r, binary.LittleEndian, &mesh.Triangles); err != nil {
		fmt.Println("Reading triangles failed: ", err)
	}

	return &mesh, nil
}

// Write writes the mesh to the given writer
func (block *Mesh) Write(w io.Writer) (int, error) {

	// return if nothing needs to be written
	if len(block.Coords) == 0 {
		return 0, nil
	}

	buf := new(bytes.Buffer)

	startCoords := meshHeaderSize
	startNormals := meshHeaderSize + len(block.Coords)*12
	startTexcoords := startNormals + len(block.Normals)*12
	startColors := startTexcoords + len(block.TexCoords)*8
	startTriangles := startColors + len(block.Colors)*12

	nameMaxLen := len(block.Name)
	if nameMaxLen > meshNameMaxSize {
		nameMaxLen = meshNameMaxSize
	}

	var data = []interface{}{
		GetDataBlockHeader(typeMesh, meshBlockVersion, block.ID, block.GetSize()),
		uint16(0), /* lod */
		uint16(0), /* maxLod */
		uint32(len(block.Coords)),
		uint32(len(block.Normals)),
		uint32(len(block.TexCoords)),
		uint32(len(block.Colors)),
		uint32(len(block.Triangles)),
		uint32(startCoords),
		uint32(startNormals),
		uint32(startTexcoords),
		uint32(startColors),
		uint32(startTriangles),
		uint64(block.MaterialID),
		uint16(len(block.Name)),
	}
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return 0, err
		}
	}

	// Name
	err := binary.Write(buf, binary.LittleEndian, []byte(block.Name[:nameMaxLen]))
	if err != nil {
		return 0, err
	}

	for i := 0; i < meshNameMaxSize-nameMaxLen; i++ {
		binary.Write(buf, binary.LittleEndian, false)
	}

	// Coords
	for _, c := range block.Coords {
		writeVec3(buf, c)
	}
	// Normals
	for _, c := range block.Normals {
		writeVec3(buf, c)
	}
	// TexCoords
	for _, c := range block.TexCoords {
		writeVec2(buf, c)
	}
	// Colors
	for _, c := range block.Colors {
		writeVec3(buf, c)
	}
	// Triangles
	for _, t := range block.Triangles {
		err := binary.Write(buf, binary.LittleEndian, t.V0)
		if err != nil {
			panic("Error during binary writing V0")
		}
		err = binary.Write(buf, binary.LittleEndian, t.V1)
		if err != nil {
			panic("Error during binary writing V1")
		}
		err = binary.Write(buf, binary.LittleEndian, t.V2)
		if err != nil {
			panic("Error during binary writing V2")
		}
	}
	return w.Write(buf.Bytes())
}

func writeVec2(w io.Writer, v mgl32.Vec2) {
	err := binary.Write(w, binary.LittleEndian, v.X())
	if err != nil {
		panic("Error during binary writing Vec2")
	}
	err = binary.Write(w, binary.LittleEndian, v.Y())
	if err != nil {
		panic("Error during binary writing Vec2")
	}
}

func writeVec3(w io.Writer, v mgl32.Vec3) {
	err := binary.Write(w, binary.LittleEndian, v.X())
	if err != nil {
		panic("Error during binary writing Vec3")
	}
	err = binary.Write(w, binary.LittleEndian, v.Y())
	if err != nil {
		panic("Error during binary writing Vec3")
	}
	err = binary.Write(w, binary.LittleEndian, v.Z())
	if err != nil {
		panic("Error during binary writing Vec3")
	}
}

// String nicely print mesh
func (m Mesh) String() string {

	s := fmt.Sprintf("\n")
	s += fmt.Sprintf("|------------------------------------------------------------|\n")
	s += fmt.Sprintf("| Mesh datablock                                             |\n")
	s += fmt.Sprintf("|------------------------------------------------------------|\n")
	s += fmt.Sprintf("| Name           | %-41s |\n", m.Name)
	s += fmt.Sprintf("| MaterialID     | %-41d |\n", m.MaterialID)
	s += fmt.Sprintf("| # Coords       | %-41d |\n", len(m.Coords))
	s += fmt.Sprintf("| # Normals      | %-41d |\n", len(m.Normals))
	s += fmt.Sprintf("| # Colors       | %-41d |\n", len(m.Colors))
	s += fmt.Sprintf("| # Triangle     | %-41d |\n", len(m.Triangles))
	s += fmt.Sprintf("\n--- Coordinates\n\n")

	for _, c := range m.Coords {
		s += fmt.Sprintf(" %5v \n", c)
	}
	s += fmt.Sprintf("\n--- Triangles\n\n")
	for _, c := range m.Triangles {
		s += fmt.Sprintf(" %5v \n", c)
	}
	s += fmt.Sprintf("|------------------------------------------------------------|\n")

	return s
}

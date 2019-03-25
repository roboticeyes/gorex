package rex

import (
	"bytes"
	"os"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestEncodingHeader(t *testing.T) {

	rexFile := File{}

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	n, err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("TEST ERROR: %v", err)
	}
	if n != 86 {
		t.Fatalf("Header size does not match")
	}
}

func TestEncodingPointList(t *testing.T) {

	t.SkipNow()

	pl := PointList{ID: 0}

	pl.Points = append(pl.Points, mgl32.Vec3{0.0, 0.0, 0.0})
	pl.Points = append(pl.Points, mgl32.Vec3{1.0, 1.0, 0.0})
	pl.Points = append(pl.Points, mgl32.Vec3{0.0, 1.0, 1.0})
	pl.Points = append(pl.Points, mgl32.Vec3{0.0, 1.0, 1.0})

	rexFile := File{}
	rexFile.PointLists = append(rexFile.PointLists, pl)

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	_, err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("TEST ERROR: %v", err)
	}

	f, _ := os.Create("pointlist.rex")
	f.Write(buf.Bytes())
	defer f.Close()
}

func TestEncodingMesh(t *testing.T) {

	// t.SkipNow()

	mesh := Mesh{ID: 1, MaterialID: 0, Name: "test"}

	mesh.Coords = append(mesh.Coords, mgl32.Vec3{0.0, 0.0, 0.0})
	mesh.Coords = append(mesh.Coords, mgl32.Vec3{1.0, 0.0, 0.0})
	mesh.Coords = append(mesh.Coords, mgl32.Vec3{0.5, 1.0, 0.0})

	mesh.Triangles = append(mesh.Triangles, Triangle{0, 1, 2})

	rexFile := File{}
	rexFile.Meshes = append(rexFile.Meshes, mesh)

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	_, err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("TEST ERROR: %v", err)
	}

	f, _ := os.Create("mesh.rex")
	f.Write(buf.Bytes())
	defer f.Close()
}

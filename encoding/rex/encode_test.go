package rex

import (
	"bytes"
	// "os"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestEncodingHeader(t *testing.T) {

	rexFile := File{}

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("TEST ERROR: %v", err)
	}
}

func TestEncodingPointList(t *testing.T) {

	pl := PointList{ID: 0}

	pl.Points = append(pl.Points, mgl32.Vec3{0.0, 0.0, 0.0})
	pl.Points = append(pl.Points, mgl32.Vec3{1.0, 1.0, 0.0})
	pl.Points = append(pl.Points, mgl32.Vec3{0.0, 1.0, 1.0})
	pl.Points = append(pl.Points, mgl32.Vec3{0.0, 1.0, 1.0})

	rexFile := File{}
	rexFile.PointLists = append(rexFile.PointLists, pl)

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("TEST ERROR: %v", err)
	}
}

func TestEncodingTrack(t *testing.T) {

	track := Track{ID: 0}

	elem := TrackElement{Point: mgl32.Vec3{0.0, 0.0, 0.0}, NormalVec: mgl32.Vec3{1.0, 1.0, 1.0}, Confidence: 1.0}
	track.Points = append(track.Points, elem)
	elem = TrackElement{Point: mgl32.Vec3{1.0, 2.0, 4.0}, NormalVec: mgl32.Vec3{1.1, 1.0, 1.0}, Confidence: 1.0}
	track.Points = append(track.Points, elem)
	elem = TrackElement{Point: mgl32.Vec3{0.7, 0.2, 2.0}, NormalVec: mgl32.Vec3{1.2, 1.0, 1.0}, Confidence: 1.0}
	track.Points = append(track.Points, elem)
	elem = TrackElement{Point: mgl32.Vec3{0.6, 0.5, 0.4}, NormalVec: mgl32.Vec3{1.3, 1.0, 1.0}, Confidence: 1.0}
	track.Points = append(track.Points, elem)

	rexFile := File{}
	rexFile.Tracks = append(rexFile.Tracks, track)

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("TEST ERROR: %v", err)
	}
}

func TestEncodingMesh(t *testing.T) {

	mesh := Mesh{ID: 1, MaterialID: 0, Name: "test"}
	mesh.Coords = append(mesh.Coords, mgl32.Vec3{0.0, 0.0, 0.0})
	mesh.Coords = append(mesh.Coords, mgl32.Vec3{1.0, 0.0, 0.0})
	mesh.Coords = append(mesh.Coords, mgl32.Vec3{0.5, 1.0, 0.0})
	mesh.Triangles = append(mesh.Triangles, Triangle{0, 1, 2})
	mesh.MaterialID = NotSpecified

	rexFile := File{}
	rexFile.Meshes = append(rexFile.Meshes, mesh)

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("TEST ERROR: %v", err)
	}
}

func TestEncodingMeshAndMaterial(t *testing.T) {

	mesh := Mesh{ID: 1, MaterialID: 0, Name: "test"}
	mesh.Coords = append(mesh.Coords, mgl32.Vec3{0.0, 0.0, 0.0})
	mesh.Coords = append(mesh.Coords, mgl32.Vec3{1.0, 0.0, 0.0})
	mesh.Coords = append(mesh.Coords, mgl32.Vec3{0.5, 1.0, 0.0})
	mesh.Triangles = append(mesh.Triangles, Triangle{0, 1, 2})
	mesh.MaterialID = 0

	mat := Material{ID: 0}
	mat.KaRgb = mgl32.Vec3{1, 0, 0}
	mat.KaTextureID = NotSpecified
	mat.KdRgb = mgl32.Vec3{1, 0, 0}
	mat.KdTextureID = NotSpecified
	mat.KsRgb = mgl32.Vec3{1, 0, 0}
	mat.KsTextureID = NotSpecified

	mat.Ns = 0
	mat.Alpha = 1

	rexFile := File{}
	rexFile.Meshes = append(rexFile.Meshes, mesh)
	rexFile.Materials = append(rexFile.Materials, mat)

	var buf bytes.Buffer
	e := NewEncoder(&buf)
	err := e.Encode(rexFile)
	if err != nil {
		t.Fatalf("TEST ERROR: %v", err)
	}

	// f, _ := os.Create("mesh.rex")
	// f.Write(buf.Bytes())
	// defer f.Close()
}

package rex

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	gorex "github.com/roboticeyes/gorex/encoding/rex"
)

// Decoder is the REX file decoder
type Decoder struct {
	r io.Reader
}

// NewDecoder opens the reader and prepares everything for building the scene graph
// If file does not exist, an error will be returned
func NewDecoder(rexFile string) (*Decoder, error) {

	file, err := os.Open(rexFile)
	if err != nil {
		return nil, fmt.Errorf("Cannot open file %s", rexFile)
	}

	r := bufio.NewReader(file)

	return &Decoder{r: r}, nil
}

// vertexNormal is used to calculate smooth normals
type vertexNormal struct {
	normals []math32.Vector3
}

// NewGroup creates and returns a group containing as children meshes.
// A group is returned even if there is only one object decoded.
func (dec *Decoder) NewGroup() (*core.Node, error) {

	group := core.NewNode()

	d := gorex.NewDecoder(dec.r)
	_, rex, err := d.Decode()

	if err != nil {
		return group, fmt.Errorf("Cannot decode REX file: %v", err)
	}
	if rex == nil {
		return group, fmt.Errorf("Nothing to decode: %v", err)
	}

	// iterate over all meshes
	for _, mesh := range rex.Meshes {

		geom := geometry.NewGeometry()
		positions := math32.NewArrayF32(0, 0)
		normals := math32.NewArrayF32(len(mesh.Coords)*3, len(mesh.Coords)*3)
		indices := math32.NewArrayU32(0, 0)
		tempNormals := make([]vertexNormal, len(mesh.Coords))

		for _, c := range mesh.Coords {
			vtx := math32.NewVector3(c.X(), c.Y(), c.Z())
			positions.AppendVector3(vtx)
		}

		for _, t := range mesh.Triangles {
			indices.Append(t.V0)
			indices.Append(t.V1)
			indices.Append(t.V2)

			// calculate normals per face
			var v0, v1, v2 math32.Vector3
			var n0, n1, n2 math32.Vector3
			var sub1, sub2 math32.Vector3
			positions.GetVector3(int(t.V0*3), &v0)
			positions.GetVector3(int(t.V1*3), &v1)
			positions.GetVector3(int(t.V2*3), &v2)
			n0.CrossVectors(sub1.SubVectors(&v1, &v0), sub2.SubVectors(&v2, &v0)).Normalize()
			n1.CrossVectors(sub1.SubVectors(&v2, &v1), sub2.SubVectors(&v0, &v1)).Normalize()
			n2.CrossVectors(sub1.SubVectors(&v0, &v2), sub2.SubVectors(&v1, &v2)).Normalize()
			tempNormals[t.V0].normals = append(tempNormals[t.V0].normals, n0)
			tempNormals[t.V1].normals = append(tempNormals[t.V1].normals, n1)
			tempNormals[t.V2].normals = append(tempNormals[t.V2].normals, n2)
		}

		// calculate smooth normals
		for i, n := range tempNormals {
			var sum math32.Vector3
			for _, normal := range n.normals {
				sum.Add(&normal)
			}
			sum.DivideScalar(float32(len(n.normals)))
			normals.SetVector3(i*3, sum.Normalize())
		}

		geom.SetIndices(indices)
		geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
		geom.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))

		// find material
		var phong *material.Phong
		for _, mat := range rex.Materials {
			if mat.ID == mesh.MaterialID {
				phong = material.NewPhong(&math32.Color{
					R: mat.KdRgb.X(),
					G: mat.KdRgb.Y(),
					B: mat.KdRgb.Z()})

				ambientColor := phong.AmbientColor()
				phong.SetAmbientColor(ambientColor.Multiply(&math32.Color{
					R: mat.KaRgb.X(),
					G: mat.KaRgb.Y(),
					B: mat.KaRgb.Z()}))
				phong.SetSpecularColor(&math32.Color{
					R: mat.KsRgb.X(),
					G: mat.KsRgb.Y(),
					B: mat.KsRgb.Z()})
				if mat.Ns != 0 {
					phong.SetShininess(mat.Ns)
				}
				break
			}
		}
		// if no material found, take a default color
		if phong == nil {
			fmt.Println("No material found, take default")
			phong = material.NewPhong(&math32.Color{R: 0.5, G: 0.2, B: 0})
		}
		// test normal shader
		// phong.SetShader("progGSDemo")

		phong.SetSide(material.SideFront)

		renderMesh := graphic.NewMesh(geom, phong)
		renderMesh.SetName(mesh.Name)
		group.Add(renderMesh)
	}

	return group, nil
}

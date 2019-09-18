package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/roboticeyes/gorex/encoding/rex"
)

func getSceneNode(id, geometryID uint64, tx, ty, tz float32) rex.SceneNode {
	return rex.SceneNode{
		ID:          id,
		GeometryID:  geometryID,
		Translation: mgl32.Vec3{tx, ty, tz},
	}
}

func main() {

	fmt.Println("Generating cube ...")

	cube, mat := rex.NewCube(1)

	rexFile := rex.File{}
	rexFile.Meshes = append(rexFile.Meshes, cube)
	rexFile.Materials = append(rexFile.Materials, mat)
	rexFile.SceneNodes = append(rexFile.SceneNodes, getSceneNode(3, 1, -5, 0, 0))
	rexFile.SceneNodes = append(rexFile.SceneNodes, getSceneNode(4, 1, 5, 0, 0))

	var buf bytes.Buffer
	e := rex.NewEncoder(&buf)
	err := e.Encode(rexFile)
	if err != nil {
		panic(err)
	}

	f, _ := os.Create("cube_mesh.rex")
	f.Write(buf.Bytes())
	defer f.Close()
}

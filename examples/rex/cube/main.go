package main

import (
	"bytes"
	"fmt"
	"math"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/roboticeyes/gorex/encoding/rex"
)

func getSceneNode(id, geometryID uint64, tx, ty, tz, scale float32) rex.SceneNode {
	return rex.SceneNode{
		ID:          id,
		GeometryID:  geometryID,
		Translation: mgl32.Vec3{tx, ty, tz},
		Scale:       mgl32.Vec3{scale, scale, scale},
	}
}

func explicit(fileName string) {
	fmt.Println("Generating cube (copy) ...")

	rexFile := rex.File{}

	var id uint64
	id = 1
	for x := -10; x < 10; x++ {
		for y := -10; y < 10; y++ {
			for z := -10; z < 10; z++ {
				cube, mat := rex.NewCube(id, id+1, 0.5)
				rexFile.Meshes = append(rexFile.Meshes, cube)
				rexFile.Materials = append(rexFile.Materials, mat)
				id += 2
			}
		}
	}
	// rexFile.SceneNodes = append(rexFile.SceneNodes, getSceneNode(4, 1, 4.5, 0, 0, 0.5))

	var buf bytes.Buffer
	e := rex.NewEncoder(&buf)
	err := e.Encode(rexFile)
	if err != nil {
		panic(err)
	}

	f, _ := os.Create(fileName)
	f.Write(buf.Bytes())
	defer f.Close()
}
func instancing(fileName string) {
	fmt.Println("Generating cube (instancing) ...")

	cube, mat := rex.NewCube(1, 2, 1)

	rexFile := rex.File{}
	rexFile.Meshes = append(rexFile.Meshes, cube)
	rexFile.Materials = append(rexFile.Materials, mat)

	id := 3
	for x := -10; x < 10; x++ {
		for y := -10; y < 10; y++ {
			for z := -10; z < 10; z++ {
				rexFile.SceneNodes = append(rexFile.SceneNodes, getSceneNode(uint64(id), 1, float32(x), float32(y), float32(z), 0.5))
				id++
			}
		}
	}
	// rexFile.SceneNodes = append(rexFile.SceneNodes, getSceneNode(4, 1, 4.5, 0, 0, 0.5))

	var buf bytes.Buffer
	e := rex.NewEncoder(&buf)
	err := e.Encode(rexFile)
	if err != nil {
		panic(err)
	}

	f, _ := os.Create(fileName)
	f.Write(buf.Bytes())
	defer f.Close()
}

func rotation(fileName string) {
	fmt.Println("Generating cube (rotation) ...")

	cube, mat := rex.NewCube(1, 2, 1)

	rexFile := rex.File{}
	rexFile.Meshes = append(rexFile.Meshes, cube)
	rexFile.Materials = append(rexFile.Materials, mat)

	rotX := FromEuler(math.Pi/4, 0, 0)
	rotY := FromEuler(0, math.Pi/4, 0)
	rotZ := FromEuler(0, 0, math.Pi/4)

	rexFile.SceneNodes = append(rexFile.SceneNodes, rex.SceneNode{
		ID:          3,
		GeometryID:  1,
		Translation: mgl32.Vec3{-5, 0, 0},
		Rotation:    mgl32.Vec4{float32(rotX.X), float32(rotX.Y), float32(rotX.Z), float32(rotX.W)},
		Scale:       mgl32.Vec3{1, 1, 1},
	})

	rexFile.SceneNodes = append(rexFile.SceneNodes, rex.SceneNode{
		ID:          4,
		GeometryID:  1,
		Translation: mgl32.Vec3{5, 0, 0},
		Rotation:    mgl32.Vec4{float32(rotY.X), float32(rotY.Y), float32(rotY.Z), float32(rotY.W)},
		Scale:       mgl32.Vec3{1, 1, 1},
	})

	rexFile.SceneNodes = append(rexFile.SceneNodes, rex.SceneNode{
		ID:          5,
		GeometryID:  1,
		Translation: mgl32.Vec3{0, 0, -5},
		Rotation:    mgl32.Vec4{float32(rotZ.X), float32(rotZ.Y), float32(rotZ.Z), float32(rotZ.W)},
		Scale:       mgl32.Vec3{1, 1, 1},
	})

	rexFile.SceneNodes = append(rexFile.SceneNodes, rex.SceneNode{
		ID:          6,
		GeometryID:  1,
		Translation: mgl32.Vec3{0, 0.5, 0},
		Scale:       mgl32.Vec3{1, 1, 1},
	})

	var buf bytes.Buffer
	e := rex.NewEncoder(&buf)
	err := e.Encode(rexFile)
	if err != nil {
		panic(err)
	}

	f, _ := os.Create(fileName)
	f.Write(buf.Bytes())
	defer f.Close()
}

func main() {

	// instancing("cube_instancing.rex")
	// explicit("cube_copy.rex")
	rotation("cube_rotated.rex")
}

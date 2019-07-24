// Copyright 2019 Robotic Eyes. All rights reserved.

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/roboticeyes/gorex/encoding/rex"
)

var (
	rexHeader  *rex.Header
	rexContent *rex.File
	// Version string from ldflags
	Version string
	// Build string from ldflags
	Build string
)

// the help text that gets displayed when something goes wrong or when you run
// help
const helpText = `
rxi - show REX file infos

actions:

  rxi -v                    prints version
  rxi help                  print this help

  rxi info "file.rex"       show all REX blocks
  rxi i "file.rex"          show all REX blocks
  rxi bbox "file.rex"       displays the bounding box of the rex file

  rxi img ID "file.rex"     extract the given image and dump it to stdout (pipe to a viewer, e.g. | feh -)
`

// help prints the help text to stdout
func help(exit int) {
	fmt.Println(helpText)
	os.Exit(exit)
}

func openRexFile(rexFile string) {
	file, err := os.Open(rexFile)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	d := rex.NewDecoder(r)
	rexHeader, rexContent, err = d.Decode()
	if err != nil && err.Error() != "unexpected EOF" {
		panic(err)
	}
}

// dumps the image to stdout (you can pipe it to an image viewer)
func rexExtractImage(rexFile, idString string) {
	openRexFile(rexFile)
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		panic(err)
	}

	for _, img := range rexContent.Images {
		if img.ID == id {
			binary.Write(os.Stdout, binary.LittleEndian, img.Data)
		}
	}
}

func rexTranslate(rexFile string, x, y, z float32, output string) {
	openRexFile(rexFile)
	fmt.Println(rexHeader)

	translate := mgl32.Vec3{x, y, z}

	if len(rexContent.Meshes) > 0 {
		for i := 0; i < len(rexContent.Meshes); i++ {
			for j := 0; j < len(rexContent.Meshes[i].Coords); j++ {
				rexContent.Meshes[i].Coords[j] = rexContent.Meshes[i].Coords[j].Add(translate)
			}
		}
	}

	// create new file
	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	e := rex.NewEncoder(&buf)
	err = e.Encode(*rexContent)
	if err != nil {
		panic(err)
	}
	n, err := f.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully written %d bytes to file %s\n", n, output)

}

func rexBbox(rexFile string) {
	openRexFile(rexFile)

	fmt.Println(rexHeader)

	bbmin := mgl32.Vec3{mgl32.MaxValue, mgl32.MaxValue, mgl32.MaxValue}
	bbmax := mgl32.Vec3{mgl32.MinValue, mgl32.MinValue, mgl32.MinValue}

	if len(rexContent.Meshes) > 0 {
		for _, mesh := range rexContent.Meshes {
			for _, c := range mesh.Coords {
				for i := 0; i < 3; i++ {
					if c[i] < bbmin[i] {
						bbmin[i] = c[i]
					}
					if c[i] > bbmax[i] {
						bbmax[i] = c[i]
					}
				}
			}
		}
	}
	fmt.Println("BoundingBox MIN: ", bbmin)
	fmt.Println("BoundingBox MAX: ", bbmax)
}

func rexInfo(rexFile string) {
	openRexFile(rexFile)

	fmt.Println(rexHeader)

	// Meshes
	if len(rexContent.Meshes) > 0 {
		fmt.Printf("Meshes (%d)\n", len(rexContent.Meshes))
		fmt.Printf("%10s %8s %8s %12s %s\n", "ID", "#Vtx", "#Tri", "Material", "Name")
		for _, mesh := range rexContent.Meshes {
			fmt.Printf("%10d %8d %8d %12d %s\n", mesh.ID, len(mesh.Coords), len(mesh.Triangles), mesh.MaterialID, mesh.Name)
		}
	}
	// Materials
	if len(rexContent.Materials) > 0 {
		fmt.Printf("Materials (%d)\n", len(rexContent.Materials))
		fmt.Printf("%10s %17s %16s %16s %5s %5s %s\n", "ID", "Ambient", "Diffuse", "Specular", "Ns", "Opacity", "TextureID (ADS)")
		for _, mat := range rexContent.Materials {
			texA, texD, texS := int(mat.KaTextureID), int(mat.KdTextureID), int(mat.KsTextureID)
			if mat.KaTextureID == rex.NotSpecified {
				texA = -1
			}
			if mat.KdTextureID == rex.NotSpecified {
				texD = -1
			}
			if mat.KsTextureID == rex.NotSpecified {
				texS = -1
			}
			fmt.Printf("%10d, [%.2f,%.2f,%.2f] [%.2f,%.2f,%.2f] [%.2f,%.2f,%.2f] %5.1f %7.2f [%d,%d,%d]\n", mat.ID,
				mat.KaRgb.X(), mat.KaRgb.Y(), mat.KaRgb.Z(),
				mat.KdRgb.X(), mat.KdRgb.Y(), mat.KdRgb.Z(),
				mat.KsRgb.X(), mat.KsRgb.Y(), mat.KsRgb.Z(),
				mat.Ns, mat.Alpha,
				texA, texD, texS)
		}
	}
	// Images
	if len(rexContent.Images) > 0 {
		fmt.Printf("Images (%d)\n", len(rexContent.Images))
		fmt.Printf("%10s %8s %12s\n", "ID", "Compression", "Bytes")
		for _, img := range rexContent.Images {
			compression := "raw"
			if img.Compression == 1 {
				compression = "jpg"
			} else if img.Compression == 2 {
				compression = "png"
			}
			fmt.Printf("%10d %11s %12d\n", img.ID, compression, len(img.Data))
		}
	}

	// PointList
	if len(rexContent.PointLists) > 0 {
		fmt.Printf("PointLists (%d)\n", len(rexContent.PointLists))
		fmt.Printf("%10s %8s %8s\n", "ID", "#Vtx", "#Col")
		for _, pl := range rexContent.PointLists {
			fmt.Printf("%10d %8d %8d\n", pl.ID, len(pl.Points), len(pl.Colors))
		}
	}

	// LineSet
	if len(rexContent.LineSets) > 0 {
		fmt.Printf("LineSets (%d)\n", len(rexContent.LineSets))
		fmt.Printf("%10s %8s %8s\n", "ID", "#Vtx", "#Col")
		for _, pl := range rexContent.LineSets {
			fmt.Printf("%10d %8d %8d\n", pl.ID, len(pl.Points), len(pl.Colors))
		}
	}

	if rexContent.UnknownBlocks > 0 {
		fmt.Printf("Unknown blocks (%d)\n", rexContent.UnknownBlocks)
	}
}

func main() {
	if len(os.Args) == 1 {
		help(0)
	}
	action := os.Args[1]

	switch action {
	case "help":
		help(0)
	case "-v":
		fmt.Printf("rxi v%s-%s\n", Version, Build)
	case "info":
		rexInfo(os.Args[2])
	case "i":
		rexInfo(os.Args[2])
	case "bbox":
		rexBbox(os.Args[2])
	case "translate":
		rexTranslate(os.Args[2], 2200, -125, 1800, "spring_infra.rex")
	case "img":
		rexExtractImage(os.Args[3], os.Args[2])
	default:
		help(1)
	}
}

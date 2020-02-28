package rex

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Decoder which can be used to read and decode REX files from a stream
type Decoder struct {
	r   io.Reader
	buf []byte
}

// NewDecoder creates a new REX decoder with a given input stream
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the input from the reader and returns
// a valid REX datastructure.
func (dec *Decoder) Decode() (*Header, *File, error) {

	header, err := ReadHeader(dec.r)
	if err != nil {
		return &Header{}, nil, err
	}
	file := &File{}

	for {
		hdr, err := ReadDataBlockHeader(dec.r)
		if err == io.EOF {
			return header, file, nil
		} else if err != nil {
			fmt.Println("*************** FOUND UNEXPECTED FILE ENDING ***************")
			return header, file, err
		}

		switch hdr.Type {
		case typeLineSet:
			ls, err := ReadLineSet(dec.r, hdr)
			if err == nil {
				file.LineSets = append(file.LineSets, *ls)
			}
		// case typeText: TODO
		case typePointList:
			pointList, err := ReadPointList(dec.r, hdr)
			if err == nil {
				file.PointLists = append(file.PointLists, *pointList)
			}
		case typeMesh:
			mesh, err := ReadMesh(dec.r, hdr)
			if err == nil {
				file.Meshes = append(file.Meshes, *mesh)
			}
		case typeImage:
			image, err := ReadImage(dec.r, hdr)
			if err == nil {
				file.Images = append(file.Images, *image)
			}
		case typeMaterial:
			material, err := ReadMaterial(dec.r, hdr)
			if err == nil {
				file.Materials = append(file.Materials, *material)
			}
		case typeSceneNode:
			sceneNode, err := ReadSceneNode(dec.r, hdr)
			if err == nil {
				file.SceneNodes = append(file.SceneNodes, *sceneNode)
			}
		case typeTrack:
			track, err := ReadTrack(dec.r, hdr)
			if err == nil {
				file.Tracks = append(file.Tracks, *track)
			}
		default:
			fmt.Printf("WARNING: Skipping type %d version %d sz %d id %d\n", hdr.Type, hdr.Version, hdr.Size, hdr.ID)
			// Read block from reader and ignore
			ignore := make([]byte, hdr.Size)
			if err := binary.Read(dec.r, binary.LittleEndian, &ignore); err != nil {
				fmt.Printf("Reading of unknown data block failed")
			}
			file.UnknownBlocks++
		}

		if err == io.EOF {
			return header, file, nil
		} else if err != nil {
			return header, nil, err
		}
	}
}

package rex

import (
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
			return header, nil, err
		}

		switch hdr.Type {
		case typeImage:
			image, err := ReadImage(dec.r, hdr)
			if err == nil {
				file.Images = append(file.Images, *image)
			}
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
		case typeMaterial:
			material, err := ReadMaterial(dec.r, hdr)
			if err == nil {
				file.Materials = append(file.Materials, *material)
			}
		default:
			fmt.Printf("WARNING: Skipping type %d version %d sz %d id %d\n", hdr.Type, hdr.Version, hdr.Size, hdr.ID)
		}

		if err == io.EOF {
			return header, file, nil
		} else if err != nil {
			return header, nil, err
		}
	}
}

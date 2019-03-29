package rex

const (
	// NotSpecified is used if no material or no texture image is set
	NotSpecified = 0x7fffffffffffffff
)

// File represents a complete valid REX file which can
// either be stored locally or sent to an arbirary writer with
// the Encoder.
type File struct {
	PointLists []PointList
	Meshes     []Mesh
	Materials  []Material
	Images     []Image
}

// Header generates a proper header for the File datastructure
func (f *File) Header() *Header {

	header := CreateHeader()

	for _, b := range f.PointLists {
		header.NrBlocks++
		header.SizeBytes += (uint64)(b.GetSize())
	}

	for _, b := range f.Meshes {
		header.NrBlocks++
		header.SizeBytes += (uint64)(b.GetSize())
	}

	for _, b := range f.Materials {
		header.NrBlocks++
		header.SizeBytes += (uint64)(b.GetSize())
	}

	for _, b := range f.Images {
		header.NrBlocks++
		header.SizeBytes += (uint64)(b.GetSize())
	}

	return header
}

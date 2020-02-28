package rex

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	trackBlockVersion = 1
)

// A TrackElement for a track consists of its x,y,z coordinate, the normalized normal vector and a
// confidence value.
type TrackElement struct {
	Point      mgl32.Vec3
	NormalVec  mgl32.Vec3
	Confidence float32
}

// Track datastructure
type Track struct {
	ID         uint64
	NrOfPoints uint32
	Timestamp  int64 // UNIX time - seconds since January 1, 1970 UTC
	Points     []TrackElement
}

// GetSize returns the estimated size of the block in bytes
func (block *Track) GetSize() int {
	return rexDataBlockHeaderSize + 4 + 8 + len(block.Points)*7
}

// ReadTrack reads a REX track
func ReadTrack(r io.Reader, hdr DataBlockHeader) (*Track, error) {
	var track Track

	if err := binary.Read(r, binary.LittleEndian, &track.NrOfPoints); err != nil {
		return nil, fmt.Errorf("Reading failed: %v ", err)
	}
	// track.NrOfPoints = nrOfPoints

	if err := binary.Read(r, binary.LittleEndian, &track.Timestamp); err != nil {
		return nil, fmt.Errorf("Reading failed: %v ", err)
	}

	track.Points = make([]TrackElement, track.NrOfPoints)
	if err := binary.Read(r, binary.LittleEndian, &track.Points); err != nil {
		return nil, fmt.Errorf("Reading coords failed: %v", err)
	}

	return &track, nil
}

// Write writes the track to the given writer
func (block *Track) Write(w io.Writer) error {
	// return if nothing needs to be written
	if len(block.Points) == 0 {
		return nil
	}

	err := WriteDataBlockHeader(w, DataBlockHeader{
		Type:    typeTrack,
		Version: trackBlockVersion,
		Size:    uint32(block.GetSize() - rexDataBlockHeaderSize),
		ID:      block.ID,
	})
	if err != nil {
		return err
	}

	var data = []interface{}{
		uint32(len(block.Points)),
		int64(block.Timestamp),
	}
	for _, v := range data {
		err := binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}
	// Track Elements
	for _, p := range block.Points {
		err := binary.Write(w, binary.LittleEndian, p.Point.X())
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, p.Point.Y())
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, p.Point.Z())
		if err != nil {
			return err
		}

		// normalize normal vector
		p.NormalVec = p.NormalVec.Normalize()

		err = binary.Write(w, binary.LittleEndian, p.NormalVec.X())
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, p.NormalVec.Y())
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, p.NormalVec.Z())
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, p.Confidence)
		if err != nil {
			return err
		}
	}

	return nil
}

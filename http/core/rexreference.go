// Copyright 2019 Robotic Eyes. All rights reserved.

package core

// RexReference is a spatial anchor which can be attached to a project or a project file.
type RexReference struct {
	Key             string                 `json:"key"`
	Project         string                 `json:"project"`
	ParentReference string                 `json:"parentReference"`
	RootReference   bool                   `json:"rootReference"`
	Address         *ProjectAddress        `json:"address"`
	AbsTransform    *ProjectTransformation `json:"absoluteTransformation"`
	RelTransform    *ProjectTransformation `json:"relativeTransformation"`
	FileTransform   *FileTransformation    `json:"fileTransformation,omitempty"`
}

// ProjectTransformation is used for the absoluteTransformation as well as for the relativeTransformation of a RexReference
type ProjectTransformation struct {
	Rotation `json:"rotation"`
	Position `json:"position"`
}

// FileTransformation is used for defining the relationship between the RexReference and the actual file.
type FileTransformation struct {
	Rotation `json:"rotation"`
	Position `json:"position"`
	Scale    float64 `json:"scale"`
}

// Position in form of a GeoJSON
type Position struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// Rotation definition given in Euler angles
type Rotation struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// NewFileTransform create a new valid FileTransformation
func NewFileTransform() *FileTransformation {
	return &FileTransformation{
		Rotation: Rotation{
			X: 0.0,
			Y: 0.0,
			Z: 0.0,
		},
		Position: Position{
			Type:        "Point",
			Coordinates: []float64{0.0, 0.0, 0.0},
		},
		Scale: 1.0,
	}
}

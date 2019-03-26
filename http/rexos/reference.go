// Copyright 2019 Robotic Eyes. All rights reserved.

package rexos

// Reference is a spatial anchor which can be attached to a project or a project file.
type Reference struct {
	Key             string                 `json:"key"`
	Project         string                 `json:"project"`
	ParentReference string                 `json:"parentReference"`
	RootReference   bool                   `json:"rootReference"`
	Address         *ProjectAddress        `json:"address"`
	AbsTransform    *ProjectTransformation `json:"absoluteTransformation"`
	RelTransform    *ProjectTransformation `json:"relativeTransformation"`
	FileTransform   *FileTransformation    `json:"fileTransformation"`
}

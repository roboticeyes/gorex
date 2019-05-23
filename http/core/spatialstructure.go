// Copyright 2019 Robotic Eyes. All rights reserved.

package core

// SpatialStructure contains the hierarchy tree of a bim model
type SpatialStructure struct {
	Name         string `json:"name"`
	GlobalID     string `json:"globalId"`
	ResourceType string `json:"resourceType"`
	ResourceUrn  string `json:"resourceUrn"`
	ResourceURL  string `json:"resourceUrl"`
	Children     []struct {
		Name         string `json:"name"`
		GlobalID     string `json:"globalId"`
		ResourceType string `json:"resourceType"`
		ResourceUrn  string `json:"resourceUrn"`
		ResourceURL  string `json:"resourceUrl"`
		Children     []struct {
			Name         string `json:"name"`
			ResourceType string `json:"resourceType"`
			ResourceUrn  string `json:"resourceUrn"`
			ResourceURL  string `json:"resourceUrl"`
			Children     []struct {
				Name         string `json:"name"`
				GlobalID     string `json:"globalId"`
				ResourceType string `json:"resourceType"`
				ResourceUrn  string `json:"resourceUrn"`
				ResourceURL  string `json:"resourceUrl"`
			} `json:"children"`
		} `json:"children"`
	} `json:"children"`
}

// Copyright 2019 Robotic Eyes. All rights reserved.

package core

// ProjectAddress defines the address information for a project
type ProjectAddress struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	AddressLine4 string `json:"addressLine4"`
	PostCode     string `json:"postcode"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Country      string `json:"country"`
}

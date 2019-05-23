// Copyright 2019 Robotic Eyes. All rights reserved.

package core

import (
	"context"
	"encoding/json"
	"strconv"
)

var (
	apiBimModels = "/bimModels"
)

// BimModelService provides the calls for accessing REX models
type BimModelService interface {
	GetBimModelByID(ctx context.Context, id uint64) (*BimModel, *SpatialStructure, HTTPStatus)
}

type bimModelService struct {
	resourceURL string // defines the URL for accessing the project resource (<schema>://<host>)
	client      *Client
}

// NewBimModelService creates a new project projectService
func NewBimModelService(client *Client, resourceURL string) BimModelService {
	return &bimModelService{
		client:      client,
		resourceURL: resourceURL,
	}
}

// GetBimModelByID returns a valid BIM model by the given ID
func (s *bimModelService) GetBimModelByID(ctx context.Context, id uint64) (*BimModel, *SpatialStructure, HTTPStatus) {

	query := s.resourceURL + apiBimModels + "/" + strconv.FormatUint(id, 10)
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &BimModel{}, &SpatialStructure{}, HTTPStatus{code, err.Error()}
	}

	var bimModel BimModel
	err = json.Unmarshal(body, &bimModel)

	spatial, status := s.getSpatialStructure(ctx, bimModel.Links.SpatialStructure.Href)
	return &bimModel, spatial, status
}

func (s *bimModelService) getSpatialStructure(ctx context.Context, url string) (*SpatialStructure, HTTPStatus) {

	body, code, err := s.client.Get(ctx, url)
	if err != nil {
		return &SpatialStructure{}, HTTPStatus{500, err.Error()}
	}

	var spatial SpatialStructure
	err = json.Unmarshal(body, &spatial)
	return &spatial, HTTPStatus{Code: code}
}

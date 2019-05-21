// Copyright 2019 Robotic Eyes. All rights reserved.

package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var (
	apiBimModels = "/bimModels"
)

// BimModelService provides the calls for accessing REX models
type BimModelService interface {
	GetBimModelByID(id uint64) (*BimModel, *SpatialStructure, HTTPStatus)
}

type bimModelService struct {
	client HTTPClient
}

// NewBimModelService creates a new project projectService
func NewBimModelService(client HTTPClient) BimModelService {
	return &bimModelService{client}
}

// GetBimModelByID returns a valid BIM model by the given ID
func (s *bimModelService) GetBimModelByID(id uint64) (*BimModel, *SpatialStructure, HTTPStatus) {

	query := s.client.GetAPIURL() + apiBimModels + "/" + strconv.FormatUint(id, 10)
	req, _ := http.NewRequest("GET", query, nil)
	body, code, err := s.client.Send(req)
	if err != nil {
		return &BimModel{}, &SpatialStructure{}, HTTPStatus{code, err.Error()}
	}

	var bimModel BimModel
	err = json.Unmarshal(body, &bimModel)

	spatial, status := s.getSpatialStructure(bimModel.Links.SpatialStructure.Href)
	return &bimModel, spatial, status
}

func (s *bimModelService) getSpatialStructure(url string) (*SpatialStructure, HTTPStatus) {

	req, _ := http.NewRequest("GET", url, nil)
	body, code, err := s.client.Send(req)
	if err != nil {
		return &SpatialStructure{}, HTTPStatus{500, err.Error()}
	}

	var spatial SpatialStructure
	err = json.Unmarshal(body, &spatial)
	return &spatial, HTTPStatus{Code: code}
}

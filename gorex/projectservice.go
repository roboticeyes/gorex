package gorex

import (
	"encoding/json"
	"net/http"
)

var (
	queryByNameAndOwner = "/api/v2/projects/search/findByNameAndOwner?"
	queryByOwner        = "/api/v2/projects/search/findAllByOwner?owner="
)

// ProjectService provides the calls for accessing REX project(s)
type ProjectService interface {
	FindAllByOwner(owner string) ([]Project, error)
	FindByNameAndOwner(name, owner string) (*Project, error)
}

type projectService struct {
	client HTTPClient
}

// NewProjectService creates a new project projectService
func NewProjectService(client HTTPClient) ProjectService {
	return &projectService{client}
}

// FindByNameAndOwner returns the unique identified project by userId and project name
func (s *projectService) FindByNameAndOwner(name, owner string) (*Project, error) {

	query := s.client.GetBaseURL() + queryByNameAndOwner + "name=" + name + "&owner=" + owner
	req, _ := http.NewRequest("GET", query, nil)
	body, err := s.client.Send(req)
	if err != nil {
		return nil, err
	}

	var project Project
	err = json.Unmarshal(body, &project)
	return &project, err
}

// FindByNameAndOwner returns the unique identified project by userId and project name
func (s *projectService) FindAllByOwner(owner string) ([]Project, error) {

	query := s.client.GetBaseURL() + queryByOwner + owner
	req, _ := http.NewRequest("GET", query, nil)
	body, err := s.client.Send(req)
	if err != nil {
		return nil, err
	}

	var projects []Project
	err = json.Unmarshal(body, &projects)
	return projects, err
}

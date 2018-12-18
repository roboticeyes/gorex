package gorex

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	queryNameAndOwner = "/api/v2/projects/search/findByNameAndOwner?"
)

// ProjectService provides the calls for accessing REX project(s)
type ProjectService interface {
	FindByNameAndOwner(name, owner string) (*Project, error)
}

type service struct {
	client HTTPClient
}

// NewProjectService creates a new project service
func NewProjectService(client HTTPClient) ProjectService {
	return &service{client}
}

func (s *service) FindByNameAndOwner(name, owner string) (*Project, error) {

	query := s.client.GetBaseURL() + queryNameAndOwner + "name=" + name + "&owner=" + owner
	req, _ := http.NewRequest("GET", query, nil)
	resp, err := s.client.Send(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
	}()

	var project Project
	err = json.NewDecoder(resp.Body).Decode(&project)
	return &project, err
}

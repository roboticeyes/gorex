// Copyright 2019 Robotic Eyes. All rights reserved.

package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

var (
	apiProjectByNameAndOwner = "/projects/search/findByNameAndOwner?"
	apiProjectByOwner        = "/projects/search/findAllByOwner?owner="
	apiProjectAllByUser      = "/projects/search/findAllFiltered?isOwnedBy=true&isReadSharedTo=true&isWriteSharedTo=true&projection=detailedList&size=100&sort=lastUpdated,desc&user="
	apiProjects              = "/projects"
	apiRexReferences         = "/rexReferences"
	apiProjectFiles          = "/projectFiles/"
)

// ProjectService provides the calls for accessing REX project(s)
type ProjectService interface {
	FindAllByUser(owner string) (*ProjectComplexList, error)
	FindAllByOwner(owner string) (*ProjectSimpleList, error)
	FindByNameAndOwner(name, owner string) (*Project, error)

	UploadProjectFile(project Project, projectFileName, fileName string, transform *FileTransformation, r io.Reader) error
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

	query := s.client.GetAPIURL() + apiProjectByNameAndOwner + "name=" + url.PathEscape(name) + "&owner=" + owner
	req, _ := http.NewRequest("GET", query, nil)
	body, err := s.client.Send(req)
	if err != nil {
		return nil, err
	}

	var project Project
	err = json.Unmarshal(body, &project)
	return &project, err
}

func (s *projectService) FindAllByUser(user string) (*ProjectComplexList, error) {
	query := s.client.GetAPIURL() + apiProjectAllByUser + user
	req, _ := http.NewRequest("GET", query, nil)
	body, err := s.client.Send(req)
	if err != nil {
		return &ProjectComplexList{}, err
	}

	var projects ProjectComplexList
	err = json.Unmarshal(body, &projects)
	if err != nil {
		panic(err)
	}

	// set ID for convenience
	for i, p := range projects.Embedded.Projects {
		re, _ := regexp.Compile("/projects/(.*)")
		values := re.FindStringSubmatch(p.Links.Self.Href)
		if len(values) > 0 {
			projects.Embedded.Projects[i].ID = values[1]
		}
	}
	return &projects, err
}

// FindByNameAndOwner returns the unique identified project by userId and project name
func (s *projectService) FindAllByOwner(owner string) (*ProjectSimpleList, error) {

	query := s.client.GetAPIURL() + apiProjectByOwner + owner
	req, _ := http.NewRequest("GET", query, nil)
	body, err := s.client.Send(req)
	if err != nil {
		return &ProjectSimpleList{}, err
	}

	var projects ProjectSimpleList
	err = json.Unmarshal(body, &projects)

	// set ID for convenience
	for i, p := range projects.Embedded.Projects {
		re, _ := regexp.Compile("/projects/(.*)")
		values := re.FindStringSubmatch(p.Links.Self.Href)
		if len(values) > 0 {
			projects.Embedded.Projects[i].ID = values[1]
		}
	}
	return &projects, err
}

// UploadProjectFile uploads a new project file.
//
// The file requires a projectFileName, which is displayed, but also a fileName which includes the suffix. The fileName
// is used for detecting the mimetype. The content of the file will be read from the `io.Reader r`.
func (s *projectService) UploadProjectFile(project Project, projectFileName, fileName string, transform *FileTransformation, r io.Reader) error {

	b := new(bytes.Buffer)

	// Get the RootRexReference
	parentReferenceURL := StripTemplateParameter(project.Embedded.RootRexReference.Links.Self.Href)

	// Create a RexReference as well
	uuid := uuid.New().String()
	rexReference := RexReference{
		Project:         project.Links.Self.Href,
		RootReference:   false,
		ParentReference: parentReferenceURL,
		Key:             uuid,
		FileTransform:   transform,
	}

	// Only create project rex reference if no one exists yet
	var rexReferenceSelfLink string
	if len(project.Embedded.RexReferences) < 2 {
		selfLink, err := s.createRexReference(&rexReference)
		if err != nil {
			return err
		}
		rexReferenceSelfLink = selfLink
	} else {
		// find non-root rex reference
		for _, r := range project.Embedded.RexReferences {
			if r.RootReference == false {
				rexReferenceSelfLink = StripTemplateParameter(r.Links.Self.Href)
				break
			}
		}
	}

	projectFile := struct {
		Name         string `json:"name"`
		Project      string `json:"project"`
		RexReference string `json:"rexReference"`
		Type         string `json:"type,omitempty"`
	}{
		Name:         projectFileName,
		Project:      project.Links.Self.Href,
		RexReference: rexReferenceSelfLink,
	}

	if filepath.Ext(fileName) == ".rex" {
		projectFile.Type = "rex"
	}

	// Create project file
	json.NewEncoder(b).Encode(projectFile)
	req, _ := http.NewRequest("POST", s.client.GetAPIURL()+apiProjectFiles, b)
	body, err := s.client.Send(req)
	if err != nil {
		return fmt.Errorf("Got server response %s with error %s", body, err)
	}

	// Upload the actual payload
	uploadURL := gjson.Get(string(body), "_links.file\\.upload.href").String()
	return s.uploadFileContent(uploadURL, fileName, r)
}

func (s *projectService) uploadFileContent(uploadURL string, fileName string, r io.Reader) error {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", fileName)
	io.Copy(part, r)
	writer.Close()

	req, _ := http.NewRequest("POST", uploadURL, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	_, err := s.client.Send(req)
	return err
}

func (s *projectService) createRexReference(r *RexReference) (string, error) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(r)

	req, _ := http.NewRequest("POST", s.client.GetAPIURL()+apiRexReferences, b)
	req.Header.Add("content-type", "application/json")
	body, err := s.client.Send(req)

	if err != nil {
		return "", fmt.Errorf("Got server response %s with error %s", body, err)
	}
	return gjson.Get(string(body), "_links.self.href").String(), nil
}

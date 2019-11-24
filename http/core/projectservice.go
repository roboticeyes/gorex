// Copyright 2019 Robotic Eyes. All rights reserved.

package core

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/roboticeyes/gorex/http/status"
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

type projectService struct {
	resourceURL string // defines the URL for accessing the project resource (<schema>://<host>)
	client      *Client
}

// NewProjectService creates a new project projectService
func NewProjectService(client *Client, resourceURL string) ProjectService {
	return &projectService{
		client:      client,
		resourceURL: resourceURL,
	}
}

// FindByNameAndOwner returns the unique identified project by userId and project name
func (s *projectService) FindByNameAndOwner(ctx context.Context, name, owner string) (*Project, status.RexReturnCode) {

	query := s.resourceURL + apiProjectByNameAndOwner + "name=" + url.PathEscape(name) + "&owner=" + owner
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &Project{}, status.RexReturnCode{Code: 500, Message: err.Error()}
	}

	var project Project
	err = json.Unmarshal(body, &project)
	return &project, status.RexReturnCode{Code: code}
}

func (s *projectService) FindAllByUser(ctx context.Context, user string, size, page uint64) (*ProjectDetailedList, status.RexReturnCode) {
	query := s.resourceURL + apiProjectAllByUser + user + "&size=" + strconv.FormatUint(size, 10) + "&page=" + strconv.FormatUint(page, 10)
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &ProjectDetailedList{}, status.RexReturnCode{Code: 500, Message: err.Error()}
	}
	if code != http.StatusOK {
		return &ProjectDetailedList{}, status.RexReturnCode{Code: code}
	}

	var projects ProjectDetailedList
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return &ProjectDetailedList{}, status.RexReturnCode{Code: 500, Message: err.Error()}
	}

	// set Urn with legacy ID if not retrieved from backend
	for i, p := range projects.Embedded.Projects {
		re, _ := regexp.Compile("/projects/(.*)")
		values := re.FindStringSubmatch(p.Links.Self.Href)
		if len(values) > 0 && p.Urn == "" {
			projects.Embedded.Projects[i].Urn = values[1]
		}
	}
	return &projects, status.RexReturnCode{Code: code}
}

// UploadProjectFile uploads a new project file.
//
// The file requires a projectFileName, which is displayed, but also a fileName which includes the suffix. The fileName
// is used for detecting the mimetype. The content of the file will be read from the `io.Reader r`.
func (s *projectService) UploadProjectFile(ctx context.Context, project Project, projectFileName, fileName string, transform *FileTransformation, r io.Reader) status.RexReturnCode {

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
		selfLink, status := s.createRexReference(ctx, &rexReference)
		if status.Code == 500 {
			return status
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
	body, code, err := s.client.Post(ctx, s.resourceURL+apiProjectFiles, b, "application/json")
	if err != nil {
		return status.RexReturnCode{Code: code, Message: err.Error()}
	}

	// Upload the actual payload
	uploadURL := gjson.Get(string(body), "_links.file\\.upload.href").String()
	return s.uploadFileContent(ctx, uploadURL, fileName, r)
}

func (s *projectService) uploadFileContent(ctx context.Context, uploadURL string, fileName string, r io.Reader) status.RexReturnCode {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", fileName)
	io.Copy(part, r)
	writer.Close()

	_, code, err := s.client.Post(ctx, uploadURL, body, writer.FormDataContentType())
	if err != nil {
		return status.RexReturnCode{Code: code, Message: err.Error()}
	}
	return status.RexReturnCode{Code: code}
}

func (s *projectService) createRexReference(ctx context.Context, r *RexReference) (string, status.RexReturnCode) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(r)

	body, code, err := s.client.Post(ctx, s.resourceURL+apiRexReferences, b, "application/json")

	if err != nil {
		return "", status.RexReturnCode{Code: 500, Message: err.Error()}
	}
	return gjson.Get(string(body), "_links.self.href").String(), status.RexReturnCode{Code: code}
}

func (s *projectService) CreateProject(ctx context.Context, name, owner string) (*Project, status.RexReturnCode) {

	project := struct {
		Name  string `json:"name"`
		Owner string `json:"owner"`
	}{
		Name:  name,
		Owner: owner,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(project)
	body, code, err := s.client.Post(ctx, s.resourceURL+apiProjects, b, "application/json")
	if err != nil {
		return nil, status.RexReturnCode{Code: code, Message: err.Error()}
	}
	if code != http.StatusCreated {
		return nil, status.RexReturnCode{Code: code, Message: "rexOS did not return 201 http code"}
	}

	var newProject Project
	err = json.Unmarshal(body, &newProject)
	if err != nil {
		return nil, status.RexReturnCode{Code: http.StatusInternalServerError, Message: "Cannot unmarshal created project"}
	}

	projectSelfLink := gjson.Get(string(body), "_links.self.href").String()
	uuid := uuid.New().String()

	// Create a RexReference as well
	rexReference := RexReference{
		Project:       projectSelfLink,
		RootReference: true,
		Key:           uuid,
	}

	_, ret := s.CreateRexReference(ctx, rexReference)
	if ret.Code != http.StatusCreated {
		// TODO delete project
		return nil, ret
	}
	return &newProject, status.RexReturnCode{Code: http.StatusCreated}
}

// CreateRexReference creates a new RexReference
func (s *projectService) CreateRexReference(ctx context.Context, r RexReference) (string, status.RexReturnCode) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(r)

	body, code, err := s.client.Post(ctx, s.resourceURL+apiRexReferences, b, "application/json")
	if err != nil {
		return "", status.RexReturnCode{Code: code, Message: err.Error()}
	}
	if code != 201 {
		return "", status.RexReturnCode{Code: code}
	}
	return gjson.Get(string(body), "_links.self.href").String(), status.RexReturnCode{Code: http.StatusCreated}
}

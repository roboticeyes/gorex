package gorex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

var (
	apiProjectByNameAndOwner = "/api/v2/projects/search/findByNameAndOwner?"
	apiProjectByOwner        = "/api/v2/projects/search/findAllByOwner?owner="
	apiProjects              = "/api/v2/projects"
	apiRexReferences         = "/api/v2/rexReferences"
	apiProjectFiles          = "/api/v2/projectFiles/"
)

// ProjectService provides the calls for accessing REX project(s)
type ProjectService interface {
	FindAllByOwner(owner string) ([]Project, error)
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

	query := s.client.GetBaseURL() + apiProjectByNameAndOwner + "name=" + name + "&owner=" + owner
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

	query := s.client.GetBaseURL() + apiProjectByOwner + owner
	req, _ := http.NewRequest("GET", query, nil)
	body, err := s.client.Send(req)
	if err != nil {
		return nil, err
	}

	var projects []Project
	err = json.Unmarshal(body, &projects)
	return projects, err
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
	rexReference := Reference{
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
	req, _ := http.NewRequest("POST", s.client.GetBaseURL()+apiProjectFiles, b)
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

func (s *projectService) createRexReference(r *Reference) (string, error) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(r)

	req, _ := http.NewRequest("POST", s.client.GetBaseURL()+apiRexReferences, b)
	req.Header.Add("content-type", "application/json")
	body, err := s.client.Send(req)

	if err != nil {
		return "", fmt.Errorf("Got server response %s with error %s", body, err)
	}
	return gjson.Get(string(body), "_links.self.href").String(), nil
}

package core

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/roboticeyes/gorex/http/creator/adding"
	"github.com/roboticeyes/gorex/http/creator/listing"
)

// Controller implements all interfaces which are required
// to work with the high-level API. It handles all the access to
// the rexOS world.
type Controller struct {
	config         RexConfig
	client         *Client
	rexUser        *User
	projectService ProjectService
	userService    UserService
}

// ProjectService provides the calls for accessing REX project(s)
type ProjectService interface {
	FindAllByUser(ctx context.Context, owner string) (*ProjectDetailedList, HTTPStatus)
	FindByNameAndOwner(ctx context.Context, name, owner string) (*Project, HTTPStatus)
	CreateProject(ctx context.Context, name, owner string) (*Project, HTTPStatus)
	UploadProjectFile(ctx context.Context, project Project, projectFileName, fileName string, transform *FileTransformation, r io.Reader) HTTPStatus
}

// UserService provides the calls for accessing REX user resource
type UserService interface {
	GetCurrentUser(ctx context.Context) (*User, HTTPStatus)
	GetTotalNumberOfUsers(ctx context.Context) (uint64, HTTPStatus)
	FindUserByEmail(ctx context.Context, email string) (*User, HTTPStatus)
	FindUserByUserID(ctx context.Context, userID string) (*User, HTTPStatus)
}

// NewController returns a new instance of the rexOS data rest interface
func NewController(config RexConfig) *Controller {
	client := NewClient()
	return &Controller{
		config:         config,
		client:         client,
		projectService: NewProjectService(client, config.ProjectResourceURL),
		userService:    NewUserService(client, config.UserResourceURL),
	}
}

// GetProjects fetches the projects and returns a list of projects
func (d *Controller) GetProjects(ctx context.Context) ([]listing.Project, error) {

	if d.rexUser == nil {
		_, err := d.getUser(ctx)
		if err != nil {
			return []listing.Project{}, err
		}
	}
	projects, status := d.projectService.FindAllByUser(ctx, d.rexUser.UserID)
	if status.Code != http.StatusOK {
		return []listing.Project{}, status
	}

	// convert data
	var result []listing.Project
	for _, p := range projects.Embedded.Projects {
		result = append(result, listing.Project{
			SelfLink:             p.Links.Self.Href,
			Urn:                  p.Urn,
			Name:                 p.Name,
			Owner:                p.Owner,
			NumberOfProjectFiles: p.NumberOfProjectFiles,
			TotalProjectFileSize: p.TotalProjectFileSize,
			Public:               p.Public,
		})
	}
	return result, nil
}

// GetUser delivers information about the authenticated user
func (d *Controller) GetUser(ctx context.Context) (listing.User, error) {

	if d.rexUser == nil {
		_, err := d.getUser(ctx)
		if err != nil {
			return listing.User{}, err
		}
	}

	return listing.User{
		UserID:    d.rexUser.UserID,
		Username:  d.rexUser.Username,
		Email:     d.rexUser.Email,
		FirstName: d.rexUser.FirstName,
		LastName:  d.rexUser.LastName,
	}, nil
}

func (d *Controller) getUser(ctx context.Context) (*User, error) {
	var status HTTPStatus
	d.rexUser, status = d.userService.GetCurrentUser(ctx)
	if status.Code != http.StatusOK {
		return nil, status
	}
	return d.rexUser, nil
}

// CreateProject create a new project with the according rex reference
func (d *Controller) CreateProject(ctx context.Context, name string) (*adding.Project, error) {
	if d.rexUser == nil {
		_, err := d.getUser(ctx)
		if err != nil {
			return nil, err
		}
	}
	p, status := d.projectService.CreateProject(ctx, name, d.rexUser.UserID)
	if status.Code != http.StatusCreated {
		return nil, status
	}
	if p == nil {
		return nil, fmt.Errorf("Did not received proper response from rexOS")
	}
	return &adding.Project{
		SelfLink: p.Links.Self.Href,
		Urn:      p.Urn,
		Name:     p.Name,
		Owner:    p.Owner,
	}, nil
}

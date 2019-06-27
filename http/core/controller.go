package core

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/roboticeyes/gorex/http/creator/adding"
	"github.com/roboticeyes/gorex/http/creator/listing"
	"github.com/roboticeyes/gorex/http/status"
)

// Controller implements all interfaces which are required
// to work with the high-level API. It handles all the access to
// the rexOS world.
type Controller struct {
	config         RexConfig
	client         *Client
	projectService ProjectService
	userService    UserService
}

// ProjectService provides the calls for accessing REX project(s)
type ProjectService interface {
	FindAllByUser(ctx context.Context, owner string, size, page uint64) (*ProjectDetailedList, status.RexReturnCode)
	FindByNameAndOwner(ctx context.Context, name, owner string) (*Project, status.RexReturnCode)
	CreateProject(ctx context.Context, name, owner string) (*Project, status.RexReturnCode)
	UploadProjectFile(ctx context.Context, project Project, projectFileName, fileName string, transform *FileTransformation, r io.Reader) status.RexReturnCode
}

// UserService provides the calls for accessing REX user resource
type UserService interface {
	GetCurrentUser(ctx context.Context) (*User, status.RexReturnCode)
	GetTotalNumberOfUsers(ctx context.Context) (uint64, status.RexReturnCode)
	FindUserByEmail(ctx context.Context, email string) (*User, status.RexReturnCode)
	FindUserByUserID(ctx context.Context, userID string) (*User, status.RexReturnCode)
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

// GetProjects fetches the projects and returns a list of projects. The size and page parameters
// are used to enable paging. The maximal number of items is limited to 100 from the rexOS backend.
func (d *Controller) GetProjects(ctx context.Context, size, page uint64) ([]listing.Project, status.RexReturnCode) {

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return []listing.Project{}, status.RexReturnCode{Code: http.StatusUnauthorized, Message: "Missing UserID in context"}
	}

	projects, retVal := d.projectService.FindAllByUser(ctx, userID, size, page)
	if retVal.Code != http.StatusOK {
		return []listing.Project{}, retVal
	}

	// convert data
	var result []listing.Project
	for _, p := range projects.Embedded.Projects {
		result = append(result, listing.Project{
			Urn:                  p.Urn,
			Name:                 p.Name,
			Owner:                p.Owner,
			NumberOfProjectFiles: p.NumberOfProjectFiles,
			TotalProjectFileSize: p.TotalProjectFileSize,
			Public:               p.Public,
		})
	}
	return result, status.RexReturnCodeOk()
}

// GetUser delivers information about the authenticated user
func (d *Controller) GetUser(ctx context.Context) (listing.User, status.RexReturnCode) {

	rexUser, retVal := d.userService.GetCurrentUser(ctx)
	if retVal.Code != http.StatusOK {
		return listing.User{}, status.RexReturnCode{Code: http.StatusUnauthorized, Message: "Cannot get user information"}
	}

	return listing.User{
		UserID:    rexUser.UserID,
		Username:  rexUser.Username,
		Email:     rexUser.Email,
		FirstName: rexUser.FirstName,
		LastName:  rexUser.LastName,
	}, status.RexReturnCodeOk()
}

// CreateProject create a new project with the according rex reference
func (d *Controller) CreateProject(ctx context.Context, name string) (*adding.Project, error) {

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Missing UserID in context")
	}
	p, status := d.projectService.CreateProject(ctx, name, userID)
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

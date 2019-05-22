package rest

import (
	"context"
	"net/http"

	"github.com/roboticeyes/gorex/http/rexos/listing"
)

// DataProviderRest provides the data from the rexOS data interface.
type DataProviderRest struct {
	client         *Client
	rexUser        *User
	projectService ProjectService
	userService    UserService
}

// NewDataProvider returns a new instance of the rexOS data rest interface
func NewDataProvider(domain string) *DataProviderRest {
	var d DataProviderRest
	d.client = NewRestClient(domain)
	d.projectService = NewProjectService(d.client)
	d.userService = NewUserService(d.client)
	return &d
}

// GetUser fetches the user information and stores it in the data provider.
func (d *DataProviderRest) GetUser(ctx context.Context) (*User, error) {
	var status HTTPStatus
	d.rexUser, status = d.userService.GetCurrentUser(ctx)
	if status.Code != http.StatusOK {
		return nil, status
	}
	return d.rexUser, nil
}

// GetProjects fetches the projects and returns a list of projects
func (d *DataProviderRest) GetProjects(ctx context.Context) ([]listing.Project, error) {

	if d.rexUser == nil {
		_, err := d.GetUser(ctx)
		if err != nil {
			return []listing.Project{}, err
		}
	}
	projects, status := d.projectService.FindAllByUser(context.Background(), d.rexUser.UserID)
	if status.Code != http.StatusOK {
		return []listing.Project{}, status
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
	return result, nil
}

// GetUserInformation delivers information about the authenticated user
func (d *DataProviderRest) GetUserInformation(ctx context.Context) (listing.User, error) {

	if d.rexUser == nil {
		_, err := d.GetUser(ctx)
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

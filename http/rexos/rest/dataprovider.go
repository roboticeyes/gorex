package rest

import (
	"context"
	"net/http"

	"github.com/roboticeyes/gorex/http/rexos/listing"
)

// DataProviderRest provides the data from the rexOS data interface. To function properly, the data provider requires
// to have a valid connected rexClient.
type DataProviderRest struct {
	rexClient      *RexClient
	client         *Client
	rexUser        *User
	projectService ProjectService
	userService    UserService
}

// NewDataProvider returns a new instance of the rexOS data rest interface
func NewDataProvider(c *RexClient) *DataProviderRest {
	var d DataProviderRest
	d.rexClient = c
	d.projectService = NewProjectService(d.client)
	d.userService = NewUserService(d.rexClient)

	var status HTTPStatus
	d.rexUser, status = d.userService.GetCurrentUser()
	if status.Code != http.StatusOK {
		panic(status)
	}
	return &d
}

// GetProjects fetches the projects and returns a list of projects
func (d *DataProviderRest) GetProjects() ([]listing.Project, error) {

	// TODO
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
func (d *DataProviderRest) GetUserInformation() (listing.User, error) {

	user, status := d.userService.GetCurrentUser()
	if status.Code != http.StatusOK {
		return listing.User{}, status
	}

	return listing.User{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

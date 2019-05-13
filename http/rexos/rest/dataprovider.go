package rest

import (
	"github.com/roboticeyes/gorex/http/rexos/listing"
)

// DataProviderRest provides the data from the rexOS data interface. To function properly, the data provider requires
// to have a valid connected rexClient.
type DataProviderRest struct {
	rexClient      *RexClient
	rexUser        *User
	projectService ProjectService
	userService    UserService
}

// NewDataProvider returns a new instance of the rexOS data rest interface
func NewDataProvider(c *RexClient) *DataProviderRest {
	var d DataProviderRest
	d.rexClient = c
	d.projectService = NewProjectService(d.rexClient)
	d.userService = NewUserService(d.rexClient)

	var err error
	d.rexUser, err = d.userService.GetCurrentUser()
	if err != nil {
		panic(err)
	}
	return &d
}

// GetProjects fetches the projects and returns a list of projects
func (d *DataProviderRest) GetProjects() ([]listing.Project, error) {

	projects, err := d.projectService.FindAllByUser(d.rexUser.UserID)
	if err != nil {
		return []listing.Project{}, err
	}

	// convert data
	var result []listing.Project
	for _, p := range projects.Embedded.Projects {
		result = append(result, listing.Project{
			ID:                   p.ID,
			Name:                 p.Name,
			Owner:                p.Owner,
			NumberOfProjectFiles: p.NumberOfProjectFiles,
			TotalProjectFileSize: p.TotalProjectFileSize,
			Public:               p.Public,
		})
	}
	return result, nil
}

package rest

import (
	"github.com/roboticeyes/gorex/http/rexos/listing"
)

// DataProviderRest provides the data from the rexOS data interface. To function properly, the data provider requires
// to have a valid connected rexClient.
type DataProviderRest struct {
	rexClient      *RexClient
	projectService ProjectService
}

// NewDataProvider returns a new instance of the rexOS data rest interface
func NewDataProvider(c *RexClient) *DataProviderRest {
	var d DataProviderRest
	d.rexClient = c
	d.projectService = NewProjectService(d.rexClient)
	return &d
}

// GetProjects fetches the projects and returns a list of projects
func (d *DataProviderRest) GetProjects() ([]listing.Project, error) {

	// TODO
	projects, err := d.projectService.FindAllByUser("initial-admin") //c.rexUser.UserID)
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

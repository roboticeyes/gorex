package listing

import (
	"context"
)

// Service providing listing functions
type Service interface {
	GetProjects(ctx context.Context) ([]Project, error)
	GetUserInformation(ctx context.Context) (User, error)
}

// DataProvider provides access to the data
type DataProvider interface {
	GetProjects(ctx context.Context) ([]Project, error)
	GetUserInformation(ctx context.Context) (User, error)
}

type listingService struct {
	provider DataProvider
}

// NewService creates an adding service with the necessary dependencies
func NewService(p DataProvider) Service {
	return &listingService{p}
}

// GetProjects returns all projects of the user
func (s *listingService) GetProjects(ctx context.Context) ([]Project, error) {
	return s.provider.GetProjects(ctx)
}

func (s *listingService) GetUserInformation(ctx context.Context) (User, error) {
	return s.provider.GetUserInformation(ctx)
}

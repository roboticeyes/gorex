package listing

import (
	"context"
)

// DataAccessor provides an abstraction of the actual data provider. Typically
// this is implemented by the core/controller to get access to the rexOS interface.
type DataAccessor interface {
	GetProjects(ctx context.Context) ([]Project, error)
	GetUser(ctx context.Context) (User, error)
}

// Service interface for getting the requested data.
type Service struct {
	dataAccessor DataAccessor
}

// NewService creates a listing service
func NewService(d DataAccessor) Service {
	return Service{d}
}

// GetProjects returns all projects of the user
func (s *Service) GetProjects(ctx context.Context) ([]Project, error) {
	return s.dataAccessor.GetProjects(ctx)
}

// GetUser gets the user information
func (s *Service) GetUser(ctx context.Context) (User, error) {
	return s.dataAccessor.GetUser(ctx)
}

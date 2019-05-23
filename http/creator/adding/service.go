package adding

import (
	"context"
)

// DataAccessor provides an abstraction of the actual data provider. Typically
// this is implemented by the core/controller to get access to the rexOS interface.
type DataAccessor interface {
	CreateProject(ctx context.Context, name string) (*Project, error)
}

// Service interface for getting the requested data.
type Service struct {
	dataAccessor DataAccessor
}

// NewService creates a listing service
func NewService(d DataAccessor) Service {
	return Service{d}
}

// CreateProject creates a new project
func (s *Service) CreateProject(ctx context.Context, name string) (*Project, error) {
	return s.dataAccessor.CreateProject(ctx, name)
}

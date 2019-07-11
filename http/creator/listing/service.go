package listing

import (
	"context"

	"github.com/roboticeyes/gorex/http/status"
)

// DataAccessor provides an abstraction of the actual data provider. Typically
// this is implemented by the core/controller to get access to the rexOS interface.
type DataAccessor interface {
	GetProjects(ctx context.Context, size, page uint64) ([]Project, status.RexReturnCode)
	GetUser(ctx context.Context) (User, status.RexReturnCode)
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
func (s *Service) GetProjects(ctx context.Context, limit, offset uint64) ([]Project, status.RexReturnCode) {
	return s.dataAccessor.GetProjects(ctx, limit, offset)
}

// GetUser gets the user information
func (s *Service) GetUser(ctx context.Context) (User, status.RexReturnCode) {
	return s.dataAccessor.GetUser(ctx)
}

package listing

// Service providing listing functions
type Service interface {
	GetProjects() ([]Project, error)
	GetUserInformation() (User, error)
}

// DataProvider provides access to the data
type DataProvider interface {
	GetProjects() ([]Project, error)
	GetUserInformation() (User, error)
}

type listingService struct {
	provider DataProvider
}

// NewService creates an adding service with the necessary dependencies
func NewService(p DataProvider) Service {
	return &listingService{p}
}

// GetProjects returns all projects of the user
func (s *listingService) GetProjects() ([]Project, error) {
	return s.provider.GetProjects()
}

func (s *listingService) GetUserInformation() (User, error) {
	return s.provider.GetUserInformation()
}

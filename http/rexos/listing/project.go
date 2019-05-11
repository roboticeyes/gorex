package listing

// Project is the top level entity of a rexOS
type Project struct {
	ID                   string
	Name                 string
	Owner                string
	NumberOfProjectFiles int
	TotalProjectFileSize int
	Public               bool
}

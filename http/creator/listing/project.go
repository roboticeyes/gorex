package listing

// Project contains the top level information for a creator project
type Project struct {
	Urn                  string
	Name                 string
	Owner                string
	NumberOfProjectFiles int
	TotalProjectFileSize int
	Public               bool
}

// ProjectFile respresents a file which can be stored to a project
type ProjectFile struct {
	Name     string
	Type     string
	FileSize int
}

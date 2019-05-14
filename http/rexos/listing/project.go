package listing

// Project is the top level entity of a rexOS
type Project struct {
	Urn                  string
	Name                 string
	Owner                string
	NumberOfProjectFiles int
	TotalProjectFileSize int
	Public               bool
}

package core

// RexConfig defines all the required settings for the rexOS DataProvider
// in order to create all rexOS service endpoints. In order to support
// internal and external usage of gorex, we need to have separate URLs
// for each resource.
type RexConfig struct {
	ProjectResourceURL string
	UserResourceURL    string
	AuthenticationURL  string
}

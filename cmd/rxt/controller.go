package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"

	"github.com/roboticeyes/gorex/http/rexos"
)

const (
	tokenFile = "token"
)

type controller struct {
	config    *Configuration
	rexClient *rexos.RexClient
	rexUser   *rexos.User
}

// NewController creates a new controller
func NewController(config *Configuration) Controller {
	return &controller{
		config: config,
	}
}

func (c *controller) Connect() (string, error) {
	return c.authenticate()
}

func (c *controller) GetConfiguration() *Configuration {
	return c.config
}

func (c *controller) GetUserID() string {
	return c.rexUser.UserID
}

func (c *controller) GetAllProjects() (rexos.ProjectComplexList, error) {
	projectService := rexos.NewProjectService(c.rexClient)
	projects, err := projectService.FindAllByUser(c.rexUser.UserID)
	if err != nil {
		return rexos.ProjectComplexList{}, err
	}
	return *projects, nil
}

// login fetches a new token and stores it locally in the token file
func (c *controller) login() error {
	cli := rexos.NewRexClient(c.config.APIUrl)

	token, err := cli.ConnectWithClientCredentials(c.config.Active.ClientID, c.config.Active.ClientSecret)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(&token)
	err = ioutil.WriteFile(tokenFile, buf, 0600)
	if err != nil {
		return err
	}
	return nil
}

// authenticate checks if a token is existing and returns a REX client
func (c *controller) authenticate() (string, error) {

	buf, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		err = c.login()
		if err != nil {
			return "", err
		}
		// try to re-read token file
		buf, err = ioutil.ReadFile(tokenFile)
		if err != nil {
			return "", err
		}
	}

	// setup client
	var token oauth2.Token
	err = json.Unmarshal(buf, &token)
	if err != nil {
		return "", err
	}
	c.rexClient = rexos.NewRexClientWithToken(c.config.APIUrl, token)

	// get user information
	userService := rexos.NewUserService(c.rexClient)

	c.rexUser, err = userService.GetCurrentUser()
	if err != nil {
		return "", err
	}
	if c.rexUser == nil || c.rexUser.UserID == "" {
		return "", fmt.Errorf("User information cannot be retrieved (token expired?), please login again")
	}
	return c.rexUser.Username, nil
}

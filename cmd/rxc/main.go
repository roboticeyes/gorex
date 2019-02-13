// Copyright 2019 Robotic Eyes. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"

	"github.com/roboticeyes/gorex"
)

// the help text that gets displayed when something goes wrong or when you run
// help
const helpText = `
rxc - command line client for rexOS

actions:

  rxc help               print this help

  rxc login              authenticate user and retrieve auth token

  rxc ls                 list all projects
  rxc ls "project name"  show details for a given project
`

const (
	tokenFile = "token"
)

var (
	apiURL       = "" // composed API url based on domain information
	scURL        = "" // composed SocketCluster url based on domain information
	clientID     = ""
	clientSecret = ""
	rexClient    *gorex.RexClient
	rexUser      *gorex.User
)

func init() {

	if os.Getenv("REX_DOMAIN") != "" {
		apiURL = "https://" + os.Getenv("REX_DOMAIN")
		scURL = "wss://" + os.Getenv("REX_DOMAIN") + "/socketcluster"
	}
	if os.Getenv("REX_CLIENT_ID") != "" {
		clientID = os.Getenv("REX_CLIENT_ID")
	}
	if os.Getenv("REX_CLIENT_SECRET") != "" {
		clientSecret = os.Getenv("REX_CLIENT_SECRET")
	}

	printSettings()
}

// help prints the help text to stdout
func help(exit int) {
	fmt.Println(helpText)
	os.Exit(exit)
}

func printSettings() {
	fmt.Println("API           ", apiURL)
	fmt.Println("SocketCluster ", scURL)
	fmt.Println("ClientID:     ", clientID)
	if clientSecret != "" {
		fmt.Println("ClientSecret:  ********")
	} else {
		fmt.Println("ClientSecret:  MISSING")
	}
}

// login fetches a new token and stores it locally in the token file
func login() {
	fmt.Println("Logging into rexOS ...")

	cli := gorex.NewRexClient(apiURL)

	token, err := cli.ConnectWithClientCredentials(clientID, clientSecret)
	if err != nil {
		log.Fatal("Error during connection", err)
	}

	buf, err := json.Marshal(&token)
	err = ioutil.WriteFile(tokenFile, buf, 0600)
	if err != nil {
		log.Fatal("Cannot write token file")
	}
	fmt.Printf("Stored token in file: %s \n\n", tokenFile)
	fmt.Println(token.AccessToken)
}

// authenticate checks if a token is existing and returns a REX client
func authenticate() {

	buf, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		login()
		// try to re-read token file
		buf, err = ioutil.ReadFile(tokenFile)
		if err != nil {
			log.Fatal("Cannot login, please contact the support team")
		}
	}

	// setup client
	var token oauth2.Token
	err = json.Unmarshal(buf, &token)
	if err != nil {
		log.Fatal("Cannot unmarshal stored token from file ", tokenFile)
	}
	rexClient = gorex.NewRexClientWithToken(apiURL, token)

	// get user information
	userService := gorex.NewUserService(rexClient)

	rexUser, err = userService.GetCurrentUser()
	if err != nil {
		log.Fatal("Cannot get user information: ", err)
	}
	if rexUser == nil || rexUser.UserID == "" {
		log.Fatal("User information cannot be retrieved (token expired?), please login again")
	}
	fmt.Println()
	fmt.Println("Logged in as")
	fmt.Println("  name:        ", rexUser.FirstName, rexUser.LastName)
	fmt.Println("  rexUsername: ", rexUser.Username)
	fmt.Println("  rexUserId:   ", rexUser.UserID)
	fmt.Println()
}

func listProjects() {
	projectService := gorex.NewProjectService(rexClient)
	projects, err := projectService.FindAllByOwner(rexUser.UserID)

	if err != nil {
		fmt.Println("Cannot get project", err)
	}

	for _, p := range projects.Embedded.Projects {
		fmt.Println("Name: ", p.Name)
		fmt.Println("Self: ", p.Links.Self.Href)
		fmt.Println()
	}
}

func listProject(projectName string) {
	projectService := gorex.NewProjectService(rexClient)
	project, err := projectService.FindByNameAndOwner(projectName, rexUser.UserID)

	if err != nil {
		fmt.Println("Cannot get project", err)
	}

	fmt.Println(project)
}

func main() {
	if len(os.Args) == 1 {
		help(0)
	}
	action := os.Args[1]
	commandArgs := len(os.Args) - 2

	switch action {
	case "help":
		help(0)
	case "login":
		login()
	case "ls":
		switch commandArgs {
		case 0:
			authenticate()
			listProjects()
		case 1:
			authenticate()
			listProject(os.Args[2])
		default:
			help(1)
		}
	default:
		help(1)
	}

	// // Setup REX client connection
	// apiURL := baseURL + "/rex-gateway/api/v2"
	// tokenURL := baseURL
	// cli := gorex.NewRexClient(tokenURL, apiURL, apiURL)
	//
	// _, err := cli.ConnectWithClientCredentials(clientID, clientSecret)
	// if err != nil {
	// 	log.Fatal("Error during connection", err)
	// }
	//
	// // Setup services
	//
	// projectService := gorex.NewProjectService(cli)
	// userService := gorex.NewUserService(cli)
	//
	// user, err := userService.GetCurrentUser()
	// if err != nil {
	// 	fmt.Println("Cannot get user", err)
	// }
	// fmt.Println(user)
	//
	// name := "test"
	// owner := user.UserID
	// project, err := projectService.FindByNameAndOwner(name, owner)
	//
	// if err != nil {
	// 	fmt.Println("Cannot get project", err)
	// }
	//
	// fmt.Println(project)
	//
	// // Add project file
	// r, _ := os.Open("/tmp/test.rex")
	// defer r.Close()
	//
	// ft := gorex.NewFileTransform()
	//
	// err = projectService.UploadProjectFile(*project, "testProjectFile", "test.rex", ft, r)
	// if err != nil {
	// 	fmt.Println("Cannot upload file", err)
	// }
}

// Copyright 2019 Robotic Eyes. All rights reserved.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/scanner"

	"github.com/breiting/socketcluster-client-go/scclient"
	"github.com/roboticeyes/gorex/http/rexos"
)

// the help text that gets displayed when something goes wrong or when you run
// help
const helpText = `
rxc - command line client for rexOS

actions:

  rxi -v                    prints version
  rxc help                  print this help

  rxc login                 authenticate user and retrieve auth token

  rxc ls                    list all projects
  rxc ls "project name"     show details for a given project
  rxc listen "project name" listens to project change notifications
`

const (
	tokenFile = "token"
)

var (
	apiURL       = "" // composed API url based on domain information
	scURL        = "" // composed SocketCluster url based on domain information
	clientID     = ""
	clientSecret = ""
	rexClient    *rexos.RexClient
	rexUser      *rexos.User
	project      *rexos.Project
	// Version string from ldflags
	Version string
	// Build string from ldflags
	Build string
)

func init() {

	if os.Getenv("REX_DOMAIN") != "" {
		apiURL = "https://" + os.Getenv("REX_DOMAIN")
		scURL = "wss://" + os.Getenv("REX_DOMAIN") + "/socketcluster/"
	}
	if os.Getenv("REX_CLIENT_ID") != "" {
		clientID = os.Getenv("REX_CLIENT_ID")
	}
	if os.Getenv("REX_CLIENT_SECRET") != "" {
		clientSecret = os.Getenv("REX_CLIENT_SECRET")
	}
}

// help prints the help text to stdout
func help(exit int) {
	fmt.Println(helpText)
	printSettings()
	os.Exit(exit)
}

func printSettings() {
	fmt.Printf("\nsettings:\n\n")
	if apiURL != "" {
		fmt.Println("  rexOS domain  ", apiURL)
	} else {
		fmt.Println("  rexOS domain   MISSING")
	}
	if scURL != "" {
		fmt.Println("  SocketCluster ", scURL)
	} else {
		fmt.Println("  SocketCluster  MISSING")
	}
	if clientID != "" {
		fmt.Println("  ClientID      ", clientID)
	} else {
		fmt.Println("  ClientID       MISSING")
	}
	if clientSecret != "" {
		fmt.Println("  ClientSecret   ********")
	} else {
		fmt.Println("  ClientSecret   MISSING")
	}
}

// login fetches a new token and stores it locally in the token file
func login() {
	fmt.Println("Logging into rexOS ...")

	cli := rexos.NewRexClient(apiURL)

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
	rexClient = rexos.NewRexClientWithToken(apiURL, token)

	// get user information
	userService := rexos.NewUserService(rexClient)

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
	projectService := rexos.NewProjectService(rexClient)
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
	projectService := rexos.NewProjectService(rexClient)
	project, err := projectService.FindByNameAndOwner(projectName, rexUser.UserID)

	if err != nil {
		fmt.Println("Cannot get project", err)
	}

	fmt.Println(project)
}

func onConnect(client scclient.Client) {
	fmt.Println("Connected to server")
}

func onDisconnect(client scclient.Client, err error) {
	fmt.Printf("Error on disconnect: %s\n", err.Error())
}

func onConnectError(client scclient.Client, err error) {
	fmt.Printf("Error on connection: %s\n", err.Error())
}

func onSetAuthentication(client scclient.Client, token string) {
	fmt.Println("Auth token received")
}

func onSocketClusterAuthentication(client scclient.Client, isAuthenticated bool) {

	// prepare proper URN for listener
	urn := "v1.resource.project." + strings.Replace(project.Urn, ":", ".", -1)

	client.SubscribeAck(urn, func(channelName string, err interface{}, data interface{}) {
		if err != nil {
			fmt.Println("Cannot get listen callback", err)
			os.Exit(1)
		}
	})

	client.OnChannel(urn, func(channelName string, data interface{}) {
		str, _ := data.(string)
		var out bytes.Buffer
		json.Indent(&out, []byte(str), "", "  ")
		out.WriteTo(os.Stdout)
	})
}

func listenProject(projectName string) {
	var reader scanner.Scanner
	var err error

	projectService := rexos.NewProjectService(rexClient)
	project, err = projectService.FindByNameAndOwner(projectName, rexUser.UserID)

	if err != nil {
		fmt.Println("Cannot get project", err)
		os.Exit(1)
	}

	fmt.Println(project)

	client := scclient.New(scURL)
	client.SetBasicListener(onConnect, onConnectError, onDisconnect)
	client.SetAuthenticationListener(onSetAuthentication, onSocketClusterAuthentication)
	client.RequestHeader = make(map[string][]string)
	authToken := "bearer " + rexClient.Token.AccessToken
	client.RequestHeader.Set("Authorization", authToken)
	go client.Connect()

	fmt.Println("Enter any key to terminate the program")
	reader.Init(os.Stdin)
	reader.Next()
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
	case "-v":
		fmt.Printf("rxi v%s-%s\n", Version, Build)
		os.Exit(0)
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
	case "listen":
		switch commandArgs {
		case 1:
			authenticate()
			listenProject(os.Args[2])
		default:
			help(1)
		}
	default:
		help(1)
	}
}

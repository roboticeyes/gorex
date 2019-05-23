// Copyright 2019 Robotic Eyes. All rights reserved.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/breiting/socketcluster-client-go/scclient"
	rexos "github.com/roboticeyes/gorex/http/core"
)

// the help text that gets displayed when something goes wrong or when you run
// help
const helpText = `
rxc - command line client for rexOS

actions:

  rxc -v                    prints version
  rxc help                  print this help

  rxc login                 authenticate user and retrieve auth token

  rxc ls                    list all projects
  rxc ls "project name"     show details for a given project
  rxc listen "project name" listens to project change notifications
  rxc bim 1000              retrieve the bim model with ID 1000
`

const (
	tokenFile = "accesstoken"
)

var (
	domain       = "" // the domain name for the rexOS (e.g. rex.robotic-eyes.com)
	apiURL       = "" // composed API url based on domain information and the api prefix
	scURL        = "" // composed SocketCluster url based on domain information
	clientID     = ""
	clientSecret = ""
	token        = "" // holds the token information after login
	rexClient    *rexos.Client
	ctx          context.Context
	project      *rexos.Project
	// Version string from ldflags
	Version string
	// Build string from ldflags
	Build string
)

func init() {

	if os.Getenv("REX_DOMAIN") != "" {
		domain = "https://" + os.Getenv("REX_DOMAIN")
		apiURL = domain + "/api/v2"
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
		fmt.Println("  rexOS API     ", apiURL)
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

	token, err := rexos.Authenticate(apiURL, clientID, clientSecret)
	if err != nil {
		log.Fatal("Error during connection", err)
	}

	err = ioutil.WriteFile(tokenFile, []byte(token), 0600)
	if err != nil {
		log.Fatal("Cannot write token file")
	}
	fmt.Printf("Stored token in file: %s \n\n", tokenFile)
	fmt.Println(token)
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

	token = string(buf)
	ctx = context.Background()
	ctx = context.WithValue(ctx, rexos.AccessTokenKey, token)
}

func listProjects() {
	client := rexos.NewClient()
	projectService := rexos.NewProjectService(client, apiURL)
	userService := rexos.NewUserService(client, apiURL)
	rexUser, status := userService.GetCurrentUser(ctx)
	if status.Code != http.StatusOK {
		fmt.Println(status)
		panic("error getting user")
	}
	projects, status := projectService.FindAllByUser(ctx, rexUser.UserID)

	if status.Code != http.StatusOK {
		fmt.Println("Cannot get project", status)
	}

	for _, p := range projects.Embedded.Projects {
		fmt.Println("Name: ", p.Name)
		fmt.Println("Self: ", p.Links.Self.Href)
		fmt.Println()
	}
}

func listProject(projectName string) {
	client := rexos.NewClient()
	projectService := rexos.NewProjectService(client, apiURL)
	userService := rexos.NewUserService(client, apiURL)
	rexUser, status := userService.GetCurrentUser(ctx)
	project, status := projectService.FindByNameAndOwner(ctx, projectName, rexUser.UserID)

	if status.Code != http.StatusOK {
		fmt.Println("Cannot get project", status)
	}

	fmt.Println(project)
}

// WIP currently not exposed, just for testing
func bimModel(modelID string) {
	bimModelService := rexos.NewBimModelService(rexClient, apiURL)
	id, err := strconv.ParseUint(modelID, 10, 64)
	if err != nil {
		id = 1000
	}
	bimModel, spatial, err := bimModelService.GetBimModelByID(ctx, id)

	if err != nil {
		fmt.Println("Cannot get project", err)
	}

	fmt.Println("BimModel:   ", bimModel.Name)
	fmt.Println("  GlobalID: ", bimModel.GlobalID)
	fmt.Println("  Owner:    ", bimModel.Owner)
	fmt.Println("  URN:      ", bimModel.Urn)
	fmt.Println()

	if spatial != nil {
		fmt.Println("Spatial:    ", spatial.Name)
		fmt.Println("  GlobalID: ", spatial.GlobalID)
		fmt.Println("  BIM site: ", spatial.Children[0].Name)

		for _, b := range spatial.Children[0].Children {
			fmt.Println("  Building: ", b.Name)
		}
	}

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

	restClient := rexos.NewClient()
	projectService := rexos.NewProjectService(restClient, apiURL)
	userService := rexos.NewUserService(restClient, apiURL)
	rexUser, status := userService.GetCurrentUser(ctx)
	if status.Code != http.StatusOK {
		fmt.Println(status)
		panic("error getting user")
	}
	project, err = projectService.FindByNameAndOwner(ctx, projectName, rexUser.UserID)

	if err != nil {
		fmt.Println("Cannot get project", err)
		os.Exit(1)
	}

	fmt.Println(project)

	client := scclient.New(scURL)
	client.SetBasicListener(onConnect, onConnectError, onDisconnect)
	client.SetAuthenticationListener(onSetAuthentication, onSocketClusterAuthentication)
	client.RequestHeader = make(map[string][]string)
	authToken := "bearer " + token
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
	case "bim":
		authenticate()
		bimModel(os.Args[2])
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

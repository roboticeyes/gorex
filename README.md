[![Build Status](https://travis-ci.org/roboticeyes/gorex.svg)](https://travis-ci.org/roboticeyes/gorex) [![Go Report Card](https://goreportcard.com/badge/github.com/roboticeyes/gorex)](https://goreportcard.com/report/github.com/roboticeyes/gorex)

# gorex

The `gorex` library provides a client implementation for the rexos API in Go. The library can easily be integrated
into your Go project. It can help you to get started with the provided REX API.

## Installation

> You can install Go by following [these instructions](https://golang.org/doc/install). Please note that Go >= 1.11. is required!

First, clone the repository to your local development path, and let go download all dependencies:

```
go mod tidy
```

This should download all required packages. To build the sample executable just use the attached `Makefile` and call
`make`.

## Usage

You can easily embed `gorex` in your target Go application by importing it as `gorex github.com/roboticeyes/gorex/gorex`


### Register an account

Visit [the REX registration](https://rex.robotic-eyes.com/registration/register) page and create a new account. Under
*Settings* you need to generate a new API token. This delivers a valid `clientId` and `clientSecret` for your
application.

### First sample

For any call into rexos you need to authenticate. Make sure that you have your `clientId` and `clientSecret` available.
You also often need your `userId` which can be found [here](https://rex.robotic-eyes.com/rex-gateway/api/v2/users/current) after
you logged into REX.

```go
	baseURL := "https://rex.robotic-eyes.com"
	clientID := "client id"
	clientSecret := "client secret"

    // Create a new client instance
	cli := gorex.NewRexClient(baseURL)

	token, err := cli.ConnectWithClientCredentials(clientID, clientSecret)
	if err != nil {
		fmt.Println("Error during connection", err)
	}

	// Create a new project service
	projectService := gorex.NewProjectService(cli)

	name := "your project name to look for"
	owner := "your user id"
	project, err := projectService.FindByNameAndOwner(name, owner)

	if err != nil {
		fmt.Println("Cannot get project", err)
	}

	fmt.Println(project)

```
## References

* [rexos](https://www.rexos.org)
* [REX](https://rex.robotic-eyes.com)

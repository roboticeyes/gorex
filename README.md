[![Build Status](https://travis-ci.org/roboticeyes/gorex.svg)](https://travis-ci.org/roboticeyes/gorex) [![Go Report Card](https://goreportcard.com/badge/github.com/roboticeyes/gorex)](https://goreportcard.com/report/github.com/roboticeyes/gorex)

<p align="center">
  <img style="float: right;" src="assets/rex-go.png" alt="goREX logo"/>
</p>

# gorex

The `gorex` library provides a library which works with [rexOS](https://www.rexos.org). The library can
easily be integrated into your Go project. It can help you to get started with REX as a developer. The library offers
two different main features:

* Working with the [REX file format](https://github.com/roboticeyes/openrex/blob/master/doc/rex-spec-v1.md)
* Working with the [rexOS REST API](https://support.robotic-eyes.com/rest/index.html)

## Installation

> You can install Go by following [these instructions](https://golang.org/doc/install). Please note that Go >= 1.11. is required!

First, clone the repository to your local development path, and let go download all dependencies:

```
go mod tidy
```

This should download all required packages. Then you can build the library by

```
go build
```

## Usage

Make sure that you just include the `gorex` library in your application:

```go
package main

import (
    "github.com/roboticeyes/gorex"
)
```

## Tools

You can easily build all tools by using the provided `Makefile`.

```
make
sudo make install
```

### rxi

`rxi` is a simple command line tool which simply dumps REX file informations to the command line. It also allows to
extract images from the file directly. For more information, please call `rxi` directly.

### rxc

`rxc` is a command line tool to work with rexOS on your command line. You can build the  `rxc` command line toole by

```go
cd cmd/rxc
go build
```

`rxc` uses environment variables to define the REX domain and user credentials, you need to set the following
environment variables:

```
REX_DOMAIN=rex.robotic-eyes.com
REX_CLIENT_ID=<your client id>
REX_CLIENT_SECRET=<your client secret>
```

Please check our [documentation](https://support.robotic-eyes.com/rest/index.html#overview-authentication) to generate
valid user credentials.

## Register an account

In order to work with the rexOS you need a REX account.
Visit [the REX registration](https://rex.robotic-eyes.com/registration/register) page and create a new account. Under
*Settings* you need to generate a new API token. This delivers a valid `clientId` and `clientSecret` for your
application.

### First sample

For any call into rexOS you need to authenticate. Make sure that you have your `clientId` and `clientSecret` available.
You also often need your `userId` which can be found [here](https://rex.robotic-eyes.com/rex-gateway/api/v2/users/current) after
you logged into REX.

```go
baseURL := "https://rex.robotic-eyes.com"
clientID := "your client id"
clientSecret := "your client secret"

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

# Todos

## REX File IO

* [ ] Data block text

## References

* [rexOS](https://www.rexos.org)
* [REX](https://rex.robotic-eyes.com)
* [REX file format v1](https://github.com/roboticeyes/openrex/blob/master/doc/rex-spec-v1.md)

# ATTENTION DEPRECATED !!!

This repository is soon deprecated. Please use our new modules:

* [GoRexFile](https://github.com/roboticeyes/gorexfile)
* [GoRexClient](https://github.com/roboticeyes/gorexclient)

If you want to use gorex in your production code, please contact support@robotic-eyes.com.

# gorex

The `gorex` library provides a library which works with [rexOS](https://www.rexos.org). The library can
easily be integrated into your Go project. It can help you to get started with REX as a developer. The library offers
two different main features:

* Working with the [REX file format](https://github.com/roboticeyes/openrex/blob/master/doc/rex-spec-v1.md)
* Working with the [rexOS REST API](https://support.robotic-eyes.com/rest/index.html)

## Architecture

The SDK provides a high-level and low-level interface to the rexOS. The `http/creator` package is a higher abstraction,
whereas the `http/core` delivers access to the low-level HAL API.

![](doc/gorex.png)

## Installation

> You can install Go by following [these instructions](https://golang.org/doc/install). Please note that Go >= 1.12. is required!

First, clone the repository to your local development path, and let go download all dependencies:

```
go mod tidy
```

This should download all required packages. To build all tools, you simple use the attached `Makefile` and call

```
make
```

## Usage

Make sure that you just include the `gorex` library in your application:

```go
package main

import (
    "github.com/roboticeyes/gorex"
)
```

Please see the `examples` folder for further demos.

## Tools

### rxi

`rxi` is a simple command line tool which simply dumps REX file informations to the command line. It also allows to
extract images from the file directly. For more information, please call `rxi` directly.

### rxt

`rxt` is terminal-based user interface for accessing the rexOS information. In order to work with `rxt`, you need to
have a configuration file in place. Either put the file into `$HOME/.config/rxt/config.json` or attach the config file
as command line parameter. The minimal information for a config file should contain the following information:

```json
{
    "default": "rex",
    "environments": [
        {
            "name": "rex",
            "domain": "rex.robotic-eyes.com",
            "clientId": "<your clientid>",
            "clientSecret": "<your clientsecret"
        }
    ]
}
```

### rxc (deprecated)

`rxc` is a command line tool to work with rexOS on your command line.

`rxc` uses environment variables to define the REX domain and user credentials, you need to set the following
environment variables:

```
REX_DOMAIN=rex.robotic-eyes.com
REX_CLIENT_ID=<your client id>
REX_CLIENT_SECRET=<your client secret>
```

Please check our [documentation](https://rexos.org) to generate valid user credentials.

## Register an account

In order to work with the rexOS you need a REX account.
Visit [the REX registration](https://rex.robotic-eyes.com/registration/register) page and create a new account. Under
*Settings* you need to generate a new API token. This delivers a valid `clientId` and `clientSecret` for your
application.

# Todos

## REX File IO

* [ ] Data block text

## References

* [rexOS](https://www.rexos.org)
* [REX](https://rex.robotic-eyes.com)
* [REX file format v1](https://github.com/roboticeyes/openrex/blob/master/doc/rex-spec-v1.md)

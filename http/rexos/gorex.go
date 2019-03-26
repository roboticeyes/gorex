// Copyright 2019 Robotic Eyes. All rights reserved.

// Package rexos provides a library for accessing the REX API. REX is a cloud-based operating system
// for building augmented reality applications.
//
// The first thing you have to do is to register at https://rex.robotic-eyes.com for a free REX
// account.  Once you activated your account, you can simply create an API access token with a
// `ClientId` and a `ClientSecret`.
package rexos

import (
	"strings"
)

// StripTemplateParameter removes the trailing template parameters of an HATEOAS URL
//
// For example: "https://rex.robotic-eyes.com/rex-gateway/api/v2/rexReferences/1000/project{?projection}"
func StripTemplateParameter(templateURL string) string {
	return strings.Split(templateURL, "{")[0]
}

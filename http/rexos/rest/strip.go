package rest

import (
	"strings"
)

// StripTemplateParameter removes the trailing template parameters of an HATEOAS URL
//
// For example: "https://rex.robotic-eyes.com/rex-gateway/api/v2/rexReferences/1000/project{?projection}"
func StripTemplateParameter(templateURL string) string {
	return strings.Split(templateURL, "{")[0]
}


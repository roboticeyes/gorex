package core

import (
	"testing"
)

func TestStripTemplateParameterOk(t *testing.T) {
	input := "https://rex.robotic-eyes.com/rex-gateway/api/v2/rexReferences/1000/project{?projection}"
	expected := "https://rex.robotic-eyes.com/rex-gateway/api/v2/rexReferences/1000/project"

	if StripTemplateParameter(input) != expected {
		t.Error("Stripping does not work")
	}
}

func TestStripTemplateParameterWithoutParameter(t *testing.T) {
	input := "https://rex.robotic-eyes.com/rex-gateway/api/v2/rexReferences/1000/project"
	expected := "https://rex.robotic-eyes.com/rex-gateway/api/v2/rexReferences/1000/project"

	if StripTemplateParameter(input) != expected {
		t.Error("Stripping does not work")
	}
}

func TestStripTemplateParameterEmpty(t *testing.T) {
	input := ""
	expected := ""

	if StripTemplateParameter(input) != expected {
		t.Error("Stripping does not work")
	}
}

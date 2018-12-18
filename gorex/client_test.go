package gorex

import (
	"golang.org/x/oauth2"
	"testing"
)

func TestCreateWithClientIdAndSecret(t *testing.T) {

	baseURL := "https://rex-test.robotic-eyes.com"

	cli := NewRexClient(baseURL)

	if cli == nil {
		t.Error("Cannot create RexClient")
	}

	// This can be used for personal testing (not part of the actual test)

	// clientID := "<clientId>"
	// clientSecret := "<clientSecret>"
	// token, err := cli.ConnectWithClientCredentials(clientID, clientSecret)
	// if token == nil {
	// 	t.Error("Token is not retrieved")
	// }
	// if err != nil {
	// 	t.Error("Error during connection")
	// }
}

func TestCreateWithToken(t *testing.T) {

	baseURL := "https://rex-test.robotic-eyes.com"
	var token oauth2.Token

	cli := NewRexClient(baseURL)
	err := cli.ConnectWithToken(token)

	if err != nil {
		t.Error("Unable to set and validate token")
	}
}

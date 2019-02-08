package gorex

import (
	"golang.org/x/oauth2"
	"testing"
)

func TestCreateWithClientIdAndSecret(t *testing.T) {

	apiURL := "https://rex-test.robotic-eyes.com/rex-gateway/api/v2"
	tokenURL := "https://rex-test.robotic-eyes.com/"

	cli := NewRexClient(tokenURL, apiURL, apiURL)

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

	apiURL := "https://rex-test.robotic-eyes.com/rex-gateway/api/v2"
	tokenURL := "https://rex-test.robotic-eyes.com/"
	var token oauth2.Token

	cli := NewRexClient(tokenURL, apiURL, apiURL)
	err := cli.ConnectWithToken(token)

	if err != nil {
		t.Error("Unable to set and validate token")
	}
}

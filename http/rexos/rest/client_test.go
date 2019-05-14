package rest

import (
	"testing"
)

func TestCreateWithClientIdAndSecret(t *testing.T) {

	apiURL := "https://rex-test.robotic-eyes.com"

	cli := NewRexClient(apiURL)

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

	apiURL := "https://rex-test.robotic-eyes.com"

	cli := NewRexClient(apiURL)
	if cli == nil {
		t.Error("Unable to set and validate token")
	}
}

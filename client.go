// Copyright 2019 Robotic Eyes. All rights reserved.

package gorex

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	apiAuth = "/oauth/token"
)

// RexClient contains the necessary data for a RexClient
type RexClient struct {
	tokenURL   string
	authURL    string
	projectURL string

	Token      oauth2.Token // Contains the authentication token
	httpClient *http.Client // The actual net client
}

// HTTPClient is an interface which is used to perform the actual
// REX request. This interface should be used for any REX API call.
// The RexClient is implementing this interface and performs the actual call.
type HTTPClient interface {
	GetTokenURL() string
	GetAuthURL() string
	GetProjectURL() string
	Send(req *http.Request) ([]byte, error)
}

// NewRexClient returns a new instance of a RexClient
func NewRexClient(tokenURL, authURL, projectURL string) *RexClient {

	return &RexClient{
		tokenURL:   tokenURL,
		authURL:    authURL,
		projectURL: projectURL,
		httpClient: http.DefaultClient,
	}
}

// NewRexClientWithToken returns a new instance of a RexClient with a given token
func NewRexClientWithToken(tokenURL, authURL, projectURL string, token oauth2.Token) *RexClient {

	return &RexClient{
		tokenURL:   tokenURL,
		authURL:    authURL,
		projectURL: projectURL,
		Token:      token,
		httpClient: http.DefaultClient,
	}
}

// ConnectWithClientCredentials performs a netowrk to the rexos backend, and retrieves
// the authentication token (stores it internally) using the given clientID and clientSecret
//
// Returns nil if connection was ok, else returns the proper error
func (c *RexClient) ConnectWithClientCredentials(clientID, clientSecret string) (*oauth2.Token, error) {

	req, err := http.NewRequest("POST", c.tokenURL+apiAuth, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return nil, err
	}

	token := clientID + ":" + clientSecret
	encodedToken := b64.StdEncoding.EncodeToString([]byte(token))
	req.Header.Add("authorization", "Basic "+encodedToken)
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=ISO-8859-1")
	req.Header.Add("accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
	}()

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Did receive HTTP code %d", resp.StatusCode)
		return nil, errors.New(msg)
	}

	err = json.Unmarshal(body, &c.Token)
	return &c.Token, err
}

// Send performs the actual HTTP call and reads the full response into a byte array which
// will be returned in case of success. Make sure that the proper token is set before making this call
func (c *RexClient) Send(req *http.Request) ([]byte, error) {

	req.Header.Add("accept", "application/json")
	c.Token.SetAuthHeader(req)
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	// this is required to properly empty the buffer for the next call
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
	}()

	return ioutil.ReadAll(resp.Body)
}

// GetProjectURL returns the REX base URL for the project resource
func (c *RexClient) GetProjectURL() string {
	return c.projectURL
}

// GetAuthURL returns the REX base URL for the authentication resource
func (c *RexClient) GetAuthURL() string {
	return c.authURL
}

// GetTokenURL returns the REX base URL for the token authentication
func (c *RexClient) GetTokenURL() string {
	return c.tokenURL
}

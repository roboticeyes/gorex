// Copyright 2019 Robotic Eyes. All rights reserved.

package rest

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

const (
	apiToken = "/oauth/token"
	apiBase  = "/rex-gateway/api/v2"
)

// RexClient contains the necessary data for a RexClient
type RexClient struct {
	domain string

	Token      oauth2.Token // Contains the authentication token
	httpClient *http.Client // The actual net client
}

// HTTPClient is an interface which is used to perform the actual
// REX request. This interface should be used for any REX API call.
// The RexClient is implementing this interface and performs the actual call.
type HTTPClient interface {
	GetAPIURL() string
	// Send performs the actual HTTP request and gets returns (body, statusCode, error)
	Send(req *http.Request) ([]byte, int, error)
}

// HTTPStatus is the return value for every REST call.If the Message is not set,
// the default status text is returned
type HTTPStatus struct {
	Code    int
	Message string
}

// Implements the error interface
func (h HTTPStatus) Error() string {
	if h.Message != "" {
		return h.Message
	}
	return http.StatusText(h.Code)
}

// NewRexClient returns a new instance of a RexClient
func NewRexClient(domain string) *RexClient {

	return &RexClient{
		domain:     domain,
		httpClient: http.DefaultClient,
	}
}

// NewRexClientWithToken returns a new instance of a RexClient with a given token
func NewRexClientWithToken(domain string, token oauth2.Token) *RexClient {

	return &RexClient{
		domain:     domain,
		Token:      token,
		httpClient: http.DefaultClient,
	}
}

// ConnectWithClientCredentials performs a netowrk to the rexos backend, and retrieves
// the authentication token (stores it internally) using the given clientID and clientSecret
//
// Returns nil if connection was ok, else returns the proper error
func (c *RexClient) ConnectWithClientCredentials(clientID, clientSecret string) (*oauth2.Token, error) {

	req, err := http.NewRequest("POST", c.getTokenURL(), strings.NewReader("grant_type=client_credentials"))
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
		msg := fmt.Sprintf("Status code %d", resp.StatusCode)
		return nil, errors.New(msg)
	}

	err = json.Unmarshal(body, &c.Token)
	return &c.Token, err
}

// Send performs the actual HTTP call and reads the full response into a byte array which
// will be returned in case of success. Make sure that the proper token is set before making this call
func (c *RexClient) Send(req *http.Request) ([]byte, int, error) {

	req.Header.Add("accept", "application/json")
	c.Token.SetAuthHeader(req)
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, resp.StatusCode, err
	}
	// this is required to properly empty the buffer for the next call
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
	}()

	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

// GetAPIURL returns the REX API URL for all API calls
func (c *RexClient) GetAPIURL() string {
	return c.domain + apiBase
}

// getTokenURL returns the REX base URL for the token authentication
func (c *RexClient) getTokenURL() string {
	return c.domain + apiToken
}

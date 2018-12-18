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
	baseURL string

	Token      oauth2.Token // Contains the authentication token
	httpClient *http.Client // The actual net client
}

// HTTPClient is an interface which is used to perform the actual
// REX request. This interface should be used for any REX API call.
// The RexClient is implementing this interface and performs the actual call.
type HTTPClient interface {
	Send(req *http.Request) (*http.Response, error)
}

// NewRexClient returns a new instance of a RexClient
func NewRexClient(baseURL string) *RexClient {

	return &RexClient{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

// NewRexClientWithToken returns a new instance of a RexClient with a given token
func NewRexClientWithToken(baseURL string, token oauth2.Token) *RexClient {

	return &RexClient{
		baseURL:    baseURL,
		Token:      token,
		httpClient: http.DefaultClient,
	}
}

// ConnectWithToken stores and validates the token for later usage
func (c *RexClient) ConnectWithToken(token oauth2.Token) error {
	c.Token = token

	// TODO validate token

	return nil
}

// ConnectWithClientCredentials performs a netowrk to the rexos backend, and retrieves
// the authentication token (stores it internally) using the given clientID and clientSecret
//
// Returns nil if connection was ok, else returns the proper error
func (c *RexClient) ConnectWithClientCredentials(clientID, clientSecret string) (*oauth2.Token, error) {

	req, err := http.NewRequest("POST", c.baseURL+apiAuth, strings.NewReader("grant_type=client_credentials"))
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

// Send fullfills the HTTPClient interface and performs a REX web request.
// Makes sure that the authentication token is available
func (c *RexClient) Send(req *http.Request) (*http.Response, error) {

	req.Header.Add("accept", "application/json")
	c.Token.SetAuthHeader(req)
	return c.httpClient.Do(req)
}

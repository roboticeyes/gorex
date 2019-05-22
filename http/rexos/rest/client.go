package rest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// AuthorizationType is used for the key of the authorization token for the context
type AuthorizationType string

const (
	// AuthorizationKey is the key for the context information
	AuthorizationKey AuthorizationType = "authorization"
)

// Client is the client which is used to send requests to the rexOS. The client
// should be created once and shared among all services.
type Client struct {
	Domain     string
	httpClient *http.Client
}

// NewRestClient create a new rexOS HTTP client
func NewRestClient(domain string) *Client {

	return &Client{
		Domain:     domain,
		httpClient: http.DefaultClient,
	}
}

// Get performs a GET request to the given query and returns the body response which is of type JSON.
// The return values also contain the http status code and a potential error which has occured.
// The request will be setup as JSON request and also takes out the authentication information from
// the given context.
func (c *Client) Get(ctx context.Context, query string) ([]byte, int, error) {

	authKey := ctx.Value(AuthorizationKey)

	if authKey == nil {
		return nil, http.StatusForbidden, fmt.Errorf("Missing token in context")
	}

	req, _ := http.NewRequest("GET", query, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("authorization", authKey.(string))
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

// Post performs a POST request to the given query, using the given payload as data, and the provided
// content-type. The content-type is typically 'application/json', but can also be of formdata in case of
// binary data upload.
func (c *Client) Post(ctx context.Context, query string, payload io.Reader, contentType string) ([]byte, int, error) {

	req, _ := http.NewRequest("POST", query, payload)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("authorization", ctx.Value("authorization").(string))
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

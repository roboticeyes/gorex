package core

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// AccessTokenType is used for the key of the authorization token for the context
type AccessTokenType string

// UserIDType is used for the key of the user_id in the context
type UserIDType string

const (
	// AccessTokenKey is the key for the context information. The context needs to store the
	// full access token with "bearer <token>"
	AccessTokenKey AccessTokenType = "authorization"
	// UserIDKey is the key for the context information. The context needs to store the
	// rexOS user id
	UserIDKey UserIDType = "UserID"
)

// Client is the client which is used to send requests to the rexOS. The client
// should be created once and shared among all services.
type Client struct {
	httpClient *http.Client
}

// NewClient create a new rexOS HTTP client
func NewClient() *Client {

	return &Client{
		httpClient: http.DefaultClient,
	}
}

// Get performs a GET request to the given query and returns the body response which is of type JSON.
// The return values also contain the http status code and a potential error which has occured.
// The request will be setup as JSON request and also takes out the authentication information from
// the given context.
func (c *Client) Get(ctx context.Context, query string) ([]byte, int, error) {

	authKey, err := GetAccessTokenFromContext(ctx)
	if err != nil {
		return nil, http.StatusForbidden, fmt.Errorf("Missing token in context")
	}

	req, _ := http.NewRequest("GET", query, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("authorization", authKey)
	resp, err := c.httpClient.Do(req)

	if err != nil {
		log.Error("GET request error:", err)
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

	authKey, err := GetAccessTokenFromContext(ctx)
	if err != nil {
		return nil, http.StatusForbidden, fmt.Errorf("Missing token in context")
	}

	req, _ := http.NewRequest("POST", query, payload)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("authorization", authKey)
	resp, err := c.httpClient.Do(req)

	if err != nil {
		log.Error("POST request error:", err)
		return nil, resp.StatusCode, err
	}
	// this is required to properly empty the buffer for the next call
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
	}()

	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

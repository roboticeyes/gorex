// Copyright 2019 Robotic Eyes. All rights reserved.

package core

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	apiToken = "/oauth/token"
)

// Authenticate gets a valid token for the given user credentials
// and returns the plain token information (without the authentication key bearer!).
// The domain denotes the rexOS base URL (e.g. <scheme>://<host>)
func Authenticate(domain, clientID, clientSecret string) (string, error) {

	client := NewClient()
	req, err := http.NewRequest("POST", domain+apiToken, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}

	token := clientID + ":" + clientSecret
	encodedToken := b64.StdEncoding.EncodeToString([]byte(token))
	req.Header.Add("authorization", "Basic "+encodedToken)
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=ISO-8859-1")
	req.Header.Add("accept", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
	}()

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Status code %d", resp.StatusCode)
	}

	var oauthToken oauth2.Token
	err = json.Unmarshal(body, &oauthToken)
	if err != nil {
		return "", err
	}
	// return "bearer <token>"
	return oauthToken.TokenType + " " + oauthToken.AccessToken, err
}

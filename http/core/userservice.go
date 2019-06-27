// Copyright 2019 Robotic Eyes. All rights reserved.

package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/roboticeyes/gorex/http/status"
	"github.com/tidwall/gjson"
)

var (
	apiCurrentUser = "/users/current"
	apiUsers       = "/users"
	apiFindByEmail = "/users/search/findUserIdByEmail?email="
	apiFindByID    = "/users/search/findByUserId?userId="
)

type userService struct {
	resourceURL string // defines the URL for accessing the project resource (<schema>://<host>)
	client      *Client
}

// NewUserService creates a new project userService
func NewUserService(client *Client, resourceURL string) UserService {
	return &userService{
		client:      client,
		resourceURL: resourceURL,
	}
}

// GetCurrentUser gets the user details of the current user.
//
// The current user is the one which has been identified by the authentication token.
func (s *userService) GetCurrentUser(ctx context.Context) (*User, status.RexReturnCode) {

	query := s.resourceURL + apiCurrentUser
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &User{}, status.RexReturnCode{code, err.Error()}
	}

	var u User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return &User{}, status.RexReturnCode{500, err.Error()}
	}
	u.SelfLink = u.Links.User.Href // assign self link
	return &u, status.RexReturnCode{Code: http.StatusOK}
}

// GetTotalNumberOfUsers returns the number of registered users.
// Requires admin permissions!
func (s *userService) GetTotalNumberOfUsers(ctx context.Context) (uint64, status.RexReturnCode) {

	query := s.resourceURL + apiUsers
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return 0, status.RexReturnCode{500, err.Error()}
	}
	return gjson.Get(string(body), "page.totalElements").Uint(), status.RexReturnCode{Code: code}
}

// FindUserByUserID returns the user information for a given user ID
// Requires admin permissions!
func (s *userService) FindUserByUserID(ctx context.Context, userID string) (*User, status.RexReturnCode) {

	query := s.resourceURL + apiFindByID + userID
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &User{}, status.RexReturnCode{500, err.Error()}
	}

	var user User
	err = json.Unmarshal(body, &user)
	return &user, status.RexReturnCode{Code: code}
}

// FindUserByEmail retrieves the user ID of a given email address
func (s *userService) FindUserByEmail(ctx context.Context, email string) (*User, status.RexReturnCode) {

	query := s.resourceURL + apiFindByEmail + email
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &User{}, status.RexReturnCode{code, err.Error()}
	}

	// check if the user can be found
	var user User
	err = json.Unmarshal(body, &user)

	if err != nil || user.UserID == "" {
		return &User{}, status.RexReturnCode{Code: http.StatusNotFound}
	}
	return &user, status.RexReturnCode{Code: http.StatusOK}
}

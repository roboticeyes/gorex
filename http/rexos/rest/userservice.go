// Copyright 2019 Robotic Eyes. All rights reserved.

package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/tidwall/gjson"
)

var (
	apiCurrentUser = "/users/current"
	apiUsers       = "/users"
	apiFindByEmail = "/users/search/findUserIdByEmail?email="
	apiFindByID    = "/users/search/findByUserId?userId="
)

// UserService provides the calls for accessing REX user resource
type UserService interface {
	GetCurrentUser(ctx context.Context) (*User, HTTPStatus)
	GetTotalNumberOfUsers(ctx context.Context) (uint64, HTTPStatus)
	FindUserByEmail(ctx context.Context, email string) (*User, HTTPStatus)
	FindUserByUserID(ctx context.Context, userID string) (*User, HTTPStatus)
}

type userService struct {
	client *Client
}

// NewUserService creates a new project userService
func NewUserService(client *Client) UserService {
	return &userService{client}
}

// GetCurrentUser gets the user details of the current user.
//
// The current user is the one which has been identified by the authentication token.
func (s *userService) GetCurrentUser(ctx context.Context) (*User, HTTPStatus) {

	query := s.client.Domain + apiCurrentUser
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &User{}, HTTPStatus{code, err.Error()}
	}

	var u User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return &User{}, HTTPStatus{500, err.Error()}
	}
	u.SelfLink = u.Links.User.Href // assign self link
	return &u, HTTPStatus{Code: http.StatusOK}
}

// GetTotalNumberOfUsers returns the number of registered users.
// Requires admin permissions!
func (s *userService) GetTotalNumberOfUsers(ctx context.Context) (uint64, HTTPStatus) {

	query := s.client.Domain + apiUsers
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return 0, HTTPStatus{500, err.Error()}
	}
	return gjson.Get(string(body), "page.totalElements").Uint(), HTTPStatus{Code: code}
}

// FindUserByUserID returns the user information for a given user ID
// Requires admin permissions!
func (s *userService) FindUserByUserID(ctx context.Context, userID string) (*User, HTTPStatus) {

	query := s.client.Domain + apiFindByID + userID
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &User{}, HTTPStatus{500, err.Error()}
	}

	var user User
	err = json.Unmarshal(body, &user)
	return &user, HTTPStatus{Code: code}
}

// FindUserByEmail retrieves the user ID of a given email address
func (s *userService) FindUserByEmail(ctx context.Context, email string) (*User, HTTPStatus) {

	query := s.client.Domain + apiFindByEmail + email
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &User{}, HTTPStatus{code, err.Error()}
	}

	// check if the user can be found
	var user User
	err = json.Unmarshal(body, &user)

	if err != nil || user.UserID == "" {
		return &User{}, HTTPStatus{Code: http.StatusNotFound}
	}
	return &user, HTTPStatus{Code: http.StatusOK}
}

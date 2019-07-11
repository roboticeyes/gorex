// Copyright 2019 Robotic Eyes. All rights reserved.

package core

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/roboticeyes/gorex/http/status"
	"github.com/tidwall/gjson"
)

var (
	apiCurrentUser = "/users/current"
	apiUsers       = "/users"
	apiFindByEmail = "/users/search/findUserIdByEmail?email="
	apiFindByID    = "/users/search/findByUserId?userId="
	apiFindAll     = "/users/?sort=lastLogin,dateCreated"
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
		return &User{}, status.RexReturnCode{Code: code, Message: err.Error()}
	}

	var u User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return &User{}, status.RexReturnCode{Code: 500, Message: err.Error()}
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
		return 0, status.RexReturnCode{Code: 500, Message: err.Error()}
	}
	return gjson.Get(string(body), "page.totalElements").Uint(), status.RexReturnCode{Code: code}
}

// FindUserByUserID returns the user information for a given user ID
// Requires admin permissions!
func (s *userService) FindUserByUserID(ctx context.Context, userID string) (*User, status.RexReturnCode) {

	query := s.resourceURL + apiFindByID + userID
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &User{}, status.RexReturnCode{Code: 500, Message: err.Error()}
	}

	var user User
	err = json.Unmarshal(body, &user)
	return &user, status.RexReturnCode{Code: code}
}

// FindUserBySelfLink returns the user based on the given self link
// Requires admin permissions!
func (s *userService) FindUserBySelfLink(ctx context.Context, selfLink string) (*User, status.RexReturnCode) {

	body, code, err := s.client.Get(ctx, selfLink)
	if err != nil {
		return &User{}, status.RexReturnCode{Code: 500, Message: err.Error()}
	}

	var user User
	err = json.Unmarshal(body, &user)
	return &user, status.RexReturnCode{Code: code}
}

// FindAllUsers returns a list of all users based on paging and size
// Requires admin permissions!
func (s *userService) FindAllUsers(ctx context.Context, size, page uint64) ([]UserDetails, status.RexReturnCode) {

	query := s.resourceURL + apiFindAll + "&page=" + strconv.FormatUint(page, 10) + "&size=" + strconv.FormatUint(size, 10)
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return []UserDetails{}, status.RexReturnCode{Code: 500, Message: err.Error()}
	}

	var userList UserList
	err = json.Unmarshal(body, &userList)
	return userList.Embedded.Users, status.RexReturnCode{Code: code}
}

func (s *userService) DeleteUser(ctx context.Context, selfLink string) status.RexReturnCode {
	err := s.client.Delete(ctx, selfLink)
	if err != nil {
		return status.RexReturnCode{Code: 500, Message: err.Error()}
	}
	return status.RexReturnCode{Code: http.StatusOK}
}

// FindUserByEmail retrieves the user ID of a given email address
func (s *userService) FindUserByEmail(ctx context.Context, email string) (*User, status.RexReturnCode) {

	query := s.resourceURL + apiFindByEmail + email
	body, code, err := s.client.Get(ctx, query)
	if err != nil {
		return &User{}, status.RexReturnCode{Code: code, Message: err.Error()}
	}

	// check if the user can be found
	var user User
	err = json.Unmarshal(body, &user)

	if err != nil || user.UserID == "" {
		return &User{}, status.RexReturnCode{Code: http.StatusNotFound}
	}
	return &user, status.RexReturnCode{Code: http.StatusOK}
}

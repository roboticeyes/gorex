// Copyright 2019 Robotic Eyes. All rights reserved.

package gorex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
)

var (
	apiCurrentUser = "/users/current"
	apiUsers       = "/users"
	apiFindByEmail = "/users/search/findUserIdByEmail?email="
	apiFindByID    = "/users/search/findByUserId?userId="
)

// UserService provides the calls for accessing REX user resource
type UserService interface {
	GetCurrentUser() (*User, error)
	GetTotalNumberOfUsers() (uint64, error)
	FindUserByEmail(email string) (*User, error)
	FindUserByUserID(userID string) (*User, error)
}

type userService struct {
	client HTTPClient
}

// NewUserService creates a new project userService
func NewUserService(client HTTPClient) UserService {
	return &userService{client}
}

// GetCurrentUser gets the user details of the current user.
//
// The current user is the one which has been identified by the authentication token.
func (s *userService) GetCurrentUser() (*User, error) {

	req, _ := http.NewRequest("GET", s.client.GetAPIURL()+apiCurrentUser, nil)

	body, err := s.client.Send(req)
	if err != nil {
		return nil, err
	}

	var u User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, fmt.Errorf("Cannot get user: %s", err)
	}
	u.SelfLink = u.Links.User.Href // assign self link
	return &u, err
}

// GetTotalNumberOfUsers returns the number of registered users.
// Requires admin permissions!
func (s *userService) GetTotalNumberOfUsers() (uint64, error) {

	req, _ := http.NewRequest("GET", s.client.GetAPIURL()+apiUsers, nil)

	body, err := s.client.Send(req)
	if err != nil {
		return 0, err
	}
	return gjson.Get(string(body), "page.totalElements").Uint(), nil
}

// FindUserByUserID returns the user information for a given user ID
// Requires admin permissions!
func (s *userService) FindUserByUserID(userID string) (*User, error) {

	req, _ := http.NewRequest("GET", s.client.GetAPIURL()+apiFindByID+userID, nil)

	body, err := s.client.Send(req)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(body, &user)
	return &user, nil
}

// FindUserByEmail retrieves the user ID of a given email address
func (s *userService) FindUserByEmail(email string) (*User, error) {

	req, _ := http.NewRequest("GET", s.client.GetAPIURL()+apiFindByEmail+email, nil)

	body, err := s.client.Send(req)
	if err != nil {
		return nil, err
	}

	// check if the user can be found
	var user User
	err = json.Unmarshal(body, &user)

	if err != nil || user.UserID == "" {
		return &User{}, errors.New("User not found")
	}
	return &user, nil
}

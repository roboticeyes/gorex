// Copyright 2018 Bernhard Reitinger. All rights reserved.

package gorex

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"net/http"
)

var (
	apiCurrentUser = "/api/v2/users/current"
	apiUsers       = "/api/v2/users"
	apiFindByEmail = "/api/v2/users/search/findUserIdByEmail?email="
	apiFindByID    = "/api/v2/users/search/findByUserId?userId="
)

// UserService provides the calls for accessing REX user resource
type UserService interface {
	GetCurrentUser() (*User, error)
	GetTotalNumberOfUsers() (uint64, error)
	GetUserByEmail(email string) (*User, error)
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

	req, _ := http.NewRequest("GET", s.client.GetBaseURL()+apiCurrentUser, nil)

	body, err := s.client.Send(req)
	if err != nil {
		return nil, err
	}

	var u User
	err = json.Unmarshal(body, &u)
	u.SelfLink = u.Links.User.Href // assign self link
	return &u, err
}

// GetTotalNumberOfUsers returns the number of registered users.
// Requires admin permissions!
func (s *userService) GetTotalNumberOfUsers() (uint64, error) {

	req, _ := http.NewRequest("GET", s.client.GetBaseURL()+apiUsers, nil)

	body, err := s.client.Send(req)
	if err != nil {
		return 0, err
	}
	return gjson.Get(string(body), "page.totalElements").Uint(), nil
}

// GetUserByEmail retrieves the user information based on a given email address
func (s *userService) GetUserByEmail(email string) (*User, error) {

	req, _ := http.NewRequest("GET", s.client.GetBaseURL()+apiFindByEmail+email, nil)

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

	// Fetch actual user information based on the retrieved UserID
	req, _ = http.NewRequest("GET", s.client.GetBaseURL()+apiFindByID+user.UserID, nil)
	body, err = s.client.Send(req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &user)
	return &user, nil
}

// Copyright 2019 Robotic Eyes. All rights reserved.

package rest

import (
	"encoding/json"
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
	GetCurrentUser() (*User, HTTPStatus)
	GetTotalNumberOfUsers() (uint64, HTTPStatus)
	FindUserByEmail(email string) (*User, HTTPStatus)
	FindUserByUserID(userID string) (*User, HTTPStatus)
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
func (s *userService) GetCurrentUser() (*User, HTTPStatus) {

	req, _ := http.NewRequest("GET", s.client.GetAPIURL()+apiCurrentUser, nil)

	body, code, err := s.client.Send(req)
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
func (s *userService) GetTotalNumberOfUsers() (uint64, HTTPStatus) {

	req, _ := http.NewRequest("GET", s.client.GetAPIURL()+apiUsers, nil)

	body, code, err := s.client.Send(req)
	if err != nil {
		return 0, HTTPStatus{500, err.Error()}
	}
	return gjson.Get(string(body), "page.totalElements").Uint(), HTTPStatus{Code: code}
}

// FindUserByUserID returns the user information for a given user ID
// Requires admin permissions!
func (s *userService) FindUserByUserID(userID string) (*User, HTTPStatus) {

	req, _ := http.NewRequest("GET", s.client.GetAPIURL()+apiFindByID+userID, nil)

	body, code, err := s.client.Send(req)
	if err != nil {
		return &User{}, HTTPStatus{500, err.Error()}
	}

	var user User
	err = json.Unmarshal(body, &user)
	return &user, HTTPStatus{Code: code}
}

// FindUserByEmail retrieves the user ID of a given email address
func (s *userService) FindUserByEmail(email string) (*User, HTTPStatus) {

	req, _ := http.NewRequest("GET", s.client.GetAPIURL()+apiFindByEmail+email, nil)

	body, code, err := s.client.Send(req)
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

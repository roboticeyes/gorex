// Copyright 2019 Robotic Eyes. All rights reserved.

package core

import (
	"fmt"
)

// User stores information of the current user.
//
// The user can either contain the currentUser information,
// but also information from a user query. The SelfLink can be
// used to directly access the data, but is also often required
// for other operations (e.g. insert a project).
type User struct {
	UserID    string `json:"userId"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	LastLogin string `json:"lastLogin,omitempty"`
	SelfLink  string
	Roles     []string `json:"roles,omitempty"`
	Links     struct {
		User struct {
			Href string `json:"href"`
		} `json:"user"`
	} `json:"_links,omitempty"`
}

// UserDetails is the full structure of a user for the /users endpoint
type UserDetails struct {
	DateCreated               string   `json:"dateCreated"`
	CreatedBy                 string   `json:"createdBy"`
	LastUpdated               string   `json:"lastUpdated"`
	UpdatedBy                 string   `json:"updatedBy"`
	UserID                    string   `json:"userId"`
	Username                  string   `json:"username"`
	Email                     string   `json:"email"`
	FirstName                 string   `json:"firstName"`
	LastName                  string   `json:"lastName"`
	ExpirationDate            string   `json:"expirationDate"`
	Locked                    bool     `json:"locked"`
	Disabled                  bool     `json:"disabled"`
	CredentialsExpirationDate string   `json:"credentialsExpirationDate"`
	LastLogin                 string   `json:"lastLogin"`
	InvitedBy                 string   `json:"invitedBy"`
	PlaintextPassword         string   `json:"plaintextPassword"`
	Roles                     []string `json:"roles"`
	Links                     struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		User struct {
			Href string `json:"href"`
		} `json:"user"`
		UserDescription struct {
			Href string `json:"href"`
		} `json:"userDescription"`
		Clients struct {
			Href string `json:"href"`
		} `json:"clients"`
		UserLicenses struct {
			Href string `json:"href"`
		} `json:"userLicenses"`
		UserConnections struct {
			Href string `json:"href"`
		} `json:"userConnections"`
		UserPaymentConnections struct {
			Href string `json:"href"`
		} `json:"userPaymentConnections"`
	} `json:"_links"`
}

// UserList is a list of users delivered by the rexos /users endpoint
type UserList struct {
	Embedded struct {
		Users []UserDetails `json:"users"`
	} `json:"_embedded"`
}

// String nicely prints out the user information.
func (u User) String() string {
	s := fmt.Sprintf("|-------------------------------------------------------------------------------|\n")
	s += fmt.Sprintf("| UserId    | %-65s |\n", u.UserID)
	s += fmt.Sprintf("| Username  | %-65s |\n", u.Username)
	s += fmt.Sprintf("| Firstname | %-65s |\n", u.FirstName)
	s += fmt.Sprintf("| Lastname  | %-65s |\n", u.LastName)
	s += fmt.Sprintf("| Email     | %-65s |\n", u.Email)
	s += fmt.Sprintf("| Self      | %-65s |\n", u.SelfLink)
	s += fmt.Sprintf("|-------------------------------------------------------------------------------|\n")

	return s
}

package gorex

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
	SelfLink  string
	Roles     []string `json:"roles,omitempty"`
	Links     struct {
		User struct {
			Href string `json:"href"`
		} `json:"user"`
	} `json:"_links,omitempty"`
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

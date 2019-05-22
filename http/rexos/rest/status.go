package rest

import (
	"net/http"
)

// HTTPStatus is the return value for every REST call.If the Message is not set,
// the default status text is returned
type HTTPStatus struct {
	Code    int
	Message string
}

// Implements the error interface
func (h HTTPStatus) Error() string {
	if h.Message != "" {
		return h.Message
	}
	return http.StatusText(h.Code)
}

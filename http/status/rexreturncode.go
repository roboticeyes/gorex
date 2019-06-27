package status

import (
	"net/http"
)

// RexReturnCode is the return value for every REST call. If the Message is not set, the default status text is
// returned. the Code reflects an http code which is returned from the REX endpoints. The status.RexReturnCode is also exposed
// in the creator package.
type RexReturnCode struct {
	Code    int
	Message string
}

// Implements the error interface
func (h RexReturnCode) Error() string {
	if h.Message != "" {
		return h.Message
	}
	return http.StatusText(h.Code)
}

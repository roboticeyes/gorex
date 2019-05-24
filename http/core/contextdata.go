package core

import (
	"context"
	"fmt"
)

const (
	// ContextDataKey is the identifier for getting the context data
	ContextDataKey = "data"
)

// ContextData is used as payload for the REX context in the interface
// The caller of the functions which takes a context must make sure that
// both data values are filled. Please use `context.WithValue` to add
// this information
type ContextData struct {
	AccessToken string
	UserID      string
}

// GetUserIDFromContext returns the user id from the context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	contextData := ctx.Value(ContextDataKey)
	if contextData == nil {
		return "", fmt.Errorf("Context does not contain any data")
	}
	return contextData.(ContextData).UserID, nil
}

// GetAccessTokenFromContext returns the accesstoken from the context
func GetAccessTokenFromContext(ctx context.Context) (string, error) {
	contextData := ctx.Value(ContextDataKey)
	if contextData == nil {
		return "", fmt.Errorf("Context does not contain any data")
	}
	return contextData.(ContextData).AccessToken, nil
}

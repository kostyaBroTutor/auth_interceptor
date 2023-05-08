package interceptor

import (
	"errors"
)

var (
	ErrVerifyToken            = errors.New("verification token error")
	ErrInvalidRequestMetadata = errors.New("invalid request metadata")
	ErrMissingAuthToken       = errors.New("missing auth token")
	ErrInvalidAuthToken       = errors.New("invalid auth token")
)

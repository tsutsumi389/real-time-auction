package service

import "errors"

// Common service errors
var (
	// Email already exists error (shared across admin and bidder services)
	ErrEmailAlreadyExists = errors.New("email already exists")
)

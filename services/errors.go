package services

import "errors"

// Custom error types for better error handling
var (
	ErrHospitalNotFound = errors.New("hospital not found")
	ErrUsernameExists   = errors.New("username already exists")
	ErrHashFailed       = errors.New("failed to hash password")
	ErrInvalidCreds     = errors.New("invalid credentials")
	ErrTokenFailed      = errors.New("failed to generate token")
)

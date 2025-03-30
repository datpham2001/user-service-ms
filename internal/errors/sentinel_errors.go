package errors

import "errors"

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrWeakPassword    = errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
)

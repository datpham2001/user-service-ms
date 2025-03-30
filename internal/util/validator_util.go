package util

import (
	"regexp"

	"github.com/datpham/user-service-ms/internal/errors"
)

var (
	emailRegex        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	minPasswordLength = 8
)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.ErrInvalidEmail
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return errors.ErrInvalidPassword
	}

	if !hasUppercase(password) || !hasLowercase(password) || !hasNumber(password) {
		return errors.ErrWeakPassword
	}

	return nil
}

func hasUppercase(password string) bool {
	return regexp.MustCompile(`[A-Z]`).MatchString(password)
}

func hasLowercase(password string) bool {
	return regexp.MustCompile(`[a-z]`).MatchString(password)
}

func hasNumber(password string) bool {
	return regexp.MustCompile(`[0-9]`).MatchString(password)
}

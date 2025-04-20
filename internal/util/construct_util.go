package util

import "fmt"

const (
	ResetPasswordTokenPrefix = "reset_password_user_id"
)

func ConstructResetPasswordTokenKey(token int) string {
	return fmt.Sprintf("%s:%d", ResetPasswordTokenPrefix, token)
}

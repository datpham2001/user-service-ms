package auth

import "context"

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Signup(ctx context.Context, req *SignupRequest) error
}

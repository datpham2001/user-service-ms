package auth

import (
	"context"

	dto "github.com/datpham/user-service-ms/internal/dto/auth"
)

type IAuthService interface {
	Signup(ctx context.Context, req *dto.UserSignupRequest) error
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)

	GetGoogleAuthUrl() string
	ProcessGoogleCallback(ctx context.Context, req *dto.GoogleCallbackRequest) (*dto.UserLoginResponse, error)
}

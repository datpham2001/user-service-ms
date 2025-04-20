package auth

import (
	"context"

	reqDto "github.com/datpham/user-service-ms/internal/dto/request"
	respDto "github.com/datpham/user-service-ms/internal/dto/response"
)

type IAuthService interface {
	Signup(ctx context.Context, req *reqDto.UserSignupRequest) error
	Login(ctx context.Context, req *reqDto.UserLoginRequest) (*respDto.UserLoginResponse, error)
	RefreshToken(ctx context.Context, req *reqDto.RefreshTokenRequest) (*respDto.UserLoginResponse, error)

	GetGoogleAuthUrl() string
	ProcessGoogleCallback(ctx context.Context, req *reqDto.GoogleCallbackRequest) (*respDto.UserLoginResponse, error)
	ForgotPassword(ctx context.Context, req *reqDto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, token int, req *reqDto.ResetPasswordRequest) error
}

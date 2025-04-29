package dto

import "github.com/datpham/user-service-ms/internal/pkg/validatorutil"

type UserSignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (r *UserSignupRequest) Validate() error {
	if err := validatorutil.ValidateEmail(r.Email); err != nil {
		return err
	}

	if err := validatorutil.ValidatePassword(r.Password); err != nil {
		return err
	}

	return nil
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type GoogleCallbackRequest struct {
	State string `form:"state" binding:"required"`
	Code  string `form:"code" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyResetPasswordTokenRequest struct {
	Token int `json:"token" binding:"required"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

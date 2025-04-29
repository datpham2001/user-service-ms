package auth

import (
	respDto "github.com/datpham/user-service-ms/internal/dto/response"
)

func (s *AuthService) mapToUserLoginResponse(accessToken string, refreshToken string) *respDto.UserLoginResponse {
	return &respDto.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

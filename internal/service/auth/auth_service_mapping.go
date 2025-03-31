package service

import dto "github.com/datpham/user-service-ms/internal/dto/auth"

func ConvertToUserLoginResponse(accessToken string, refreshToken string) *dto.UserLoginResponse {
	return &dto.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

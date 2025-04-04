package auth

import dto "github.com/datpham/user-service-ms/internal/dto/auth"

func (s *AuthService) mapToUserLoginResponse(accessToken string, refreshToken string) *dto.UserLoginResponse {
	return &dto.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

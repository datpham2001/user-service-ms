package service

import (
	"context"
	"errors"

	dto "github.com/datpham/user-service-ms/internal/dto/auth"
	customErr "github.com/datpham/user-service-ms/internal/errors"
	"github.com/datpham/user-service-ms/internal/repository/entity"
	"github.com/datpham/user-service-ms/internal/util"
	"github.com/datpham/user-service-ms/pkg/logger"
	"gorm.io/gorm"
)

type AuthService struct {
	logger         *logger.Logger
	authRepository IAuthRepository
	jwtToken       IJwtToken
}

func New(logger *logger.Logger, authRepository IAuthRepository, jwtToken IJwtToken) *AuthService {
	return &AuthService{
		logger:         logger,
		authRepository: authRepository,
		jwtToken:       jwtToken,
	}
}

// func (s *AuthService) GetAuthURL() string {
// 	return s.config.AuthCodeURL("state")
// }

// func (s *AuthService) GetUserInfo(code string) (*GoogleUserInfo, error) {
// 	token, err := s.config.Exchange(context.Background(), code)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to exchange token: %v", err)
// 	}

// 	client := s.config.Client(context.Background(), token)
// 	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get user info: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	data, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response: %v", err)
// 	}

// 	var userInfo GoogleUserInfo
// 	if err := json.Unmarshal(data, &userInfo); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
// 	}

// 	return &userInfo, nil
// }

func (s *AuthService) Signup(ctx context.Context, req *dto.UserSignupRequest) error {
	user, err := s.authRepository.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if user != nil {
		return &customErr.CustomError{
			Code:    customErr.ErrInvalidRequest,
			Message: "Email already exists",
		}
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return &customErr.CustomError{
			Code:    customErr.ErrInvalidRequest,
			Message: "Failed to hash password",
		}
	}

	user = &entity.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.authRepository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	user, err := s.authRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customErr.CustomError{
				Code:    customErr.ErrNotFound,
				Message: "User not found",
			}
		}

		return nil, err
	}

	if err := util.CheckPassword(user.Password, req.Password); err != nil {
		return nil, &customErr.CustomError{
			Code:    customErr.ErrInvalidRequest,
			Message: "Incorrect password",
		}
	}

	accessToken, refreshToken, err := s.jwtToken.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, err
	}

	return ConvertToUserLoginResponse(accessToken, refreshToken), nil
}

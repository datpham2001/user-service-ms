package auth

import (
	"context"
	"errors"
	"fmt"

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
	jwtTokenSvc    IJwtTokenService
	oauthSvc       IOAuthService
}

func New(
	logger *logger.Logger,
	authRepository IAuthRepository,
	jwtTokenSvc IJwtTokenService,
	oauthSvc IOAuthService,
) *AuthService {
	return &AuthService{
		logger:         logger,
		authRepository: authRepository,
		jwtTokenSvc:    jwtTokenSvc,
		oauthSvc:       oauthSvc,
	}
}

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
		return err
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

	accessToken, refreshToken, err := s.jwtTokenSvc.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.authRepository.UpdateById(ctx, user.ID, &entity.User{
		RefreshToken: refreshToken,
	}); err != nil {
		return nil, err
	}

	return s.mapToUserLoginResponse(accessToken, refreshToken), nil
}

func (s *AuthService) GetGoogleAuthUrl() string {
	return s.oauthSvc.GetGoogleAuthUrl()
}

func (s *AuthService) ProcessGoogleCallback(ctx context.Context, req *dto.GoogleCallbackRequest) (*dto.UserLoginResponse, error) {
	if err := s.oauthSvc.VerifyGoogleState(req.State); err != nil {
		return nil, fmt.Errorf("failed to verify google state: %w", err)
	}

	token, err := s.oauthSvc.GetGoogleAccessToken(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to get google access token: %w", err)
	}

	userInfo, err := s.oauthSvc.GetGoogleUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get google user info: %w", err)
	}

	user, err := s.authRepository.GetByEmail(ctx, userInfo["email"].(string))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user != nil {
		return nil, &customErr.CustomError{
			Code:    customErr.ErrInvalidRequest,
			Message: "Email already exists",
		}
	}

	user = &entity.User{
		Email: userInfo["email"].(string),
	}

	if err := s.authRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	return nil, nil
}

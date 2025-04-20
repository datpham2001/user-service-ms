package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	reqDto "github.com/datpham/user-service-ms/internal/dto/request"
	respDto "github.com/datpham/user-service-ms/internal/dto/response"
	customErr "github.com/datpham/user-service-ms/internal/errors"
	"github.com/datpham/user-service-ms/internal/infra/rabbitmq"
	"github.com/datpham/user-service-ms/internal/repository/entity"
	"github.com/datpham/user-service-ms/internal/util"
	"github.com/datpham/user-service-ms/pkg/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	ResetPasswordTokenTTL = time.Minute * 10
)

type AuthService struct {
	logger         *logger.Logger
	authRepository IAuthRepository
	jwtTokenSvc    IJwtTokenService
	oauthSvc       IOAuthService
	cacheSvc       ICacheService
	rabbitMQ       *rabbitmq.RabbitMQ
}

func New(
	logger *logger.Logger,
	authRepository IAuthRepository,
	jwtTokenSvc IJwtTokenService,
	oauthSvc IOAuthService,
	cacheSvc ICacheService,
	rabbitMQ *rabbitmq.RabbitMQ,
) *AuthService {
	return &AuthService{
		logger:         logger,
		authRepository: authRepository,
		jwtTokenSvc:    jwtTokenSvc,
		oauthSvc:       oauthSvc,
		cacheSvc:       cacheSvc,
		rabbitMQ:       rabbitMQ,
	}
}

func (s *AuthService) Signup(ctx context.Context, req *reqDto.UserSignupRequest) error {
	user, err := s.authRepository.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get user by email: %w", err)
	}
	if user != nil {
		return customErr.NewCustomError(customErr.ErrInvalidRequest, "Email already exists")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user = &entity.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.authRepository.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, req *reqDto.UserLoginRequest) (*respDto.UserLoginResponse, error) {
	user, err := s.authRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customErr.NewCustomError(customErr.ErrNotFound, "User not found")
		}

		return nil, fmt.Errorf("failed to get user by email: %s", err.Error())
	}

	if err := util.CheckPassword(user.Password, req.Password); err != nil {
		return nil, customErr.NewCustomError(customErr.ErrInvalidRequest, "Incorrect password")
	}

	accessToken, refreshToken, err := s.jwtTokenSvc.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %s", err.Error())
	}

	if err := s.authRepository.UpdateById(ctx, user.ID, &entity.User{
		RefreshToken: refreshToken,
	}); err != nil {
		return nil, fmt.Errorf("failed to update user: %s", err.Error())
	}

	return s.mapToUserLoginResponse(accessToken, refreshToken), nil
}

func (s *AuthService) GetGoogleAuthUrl() string {
	return s.oauthSvc.GetGoogleAuthUrl()
}

func (s *AuthService) ProcessGoogleCallback(
	ctx context.Context,
	req *reqDto.GoogleCallbackRequest,
) (*respDto.UserLoginResponse, error) {
	if err := s.oauthSvc.VerifyGoogleState(req.State); err != nil {
		return nil, fmt.Errorf("failed to verify google state: %s", err.Error())
	}

	token, err := s.oauthSvc.GetGoogleAccessToken(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to get google access token: %s", err.Error())
	}

	userInfo, err := s.oauthSvc.GetGoogleUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get google user info: %s", err.Error())
	}

	user, err := s.authRepository.GetByEmail(ctx, userInfo["email"].(string))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user != nil {
		return nil, customErr.NewCustomError(customErr.ErrInvalidRequest, "Email already exists")
	}

	user = &entity.User{
		Email: userInfo["email"].(string),
	}

	if err := s.authRepository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %s", err.Error())
	}

	return nil, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *reqDto.RefreshTokenRequest) (*respDto.UserLoginResponse, error) {
	user, err := s.authRepository.GetByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, customErr.NewCustomError(customErr.ErrNotFound, "User not found")
	}

	accessToken, refreshToken, err := s.jwtTokenSvc.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %s", err.Error())
	}

	if err := s.authRepository.UpdateById(ctx, user.ID, &entity.User{
		RefreshToken: refreshToken,
	}); err != nil {
		return nil, err
	}

	return s.mapToUserLoginResponse(accessToken, refreshToken), nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req *reqDto.ForgotPasswordRequest) error {
	user, err := s.authRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customErr.NewCustomError(customErr.ErrNotFound, "User not found")
		}

		return fmt.Errorf("failed to get user by email: %s", err.Error())
	}

	if user == nil {
		return customErr.NewCustomError(customErr.ErrNotFound, "User not found")
	}

	token, err := util.GenerateResetPasswordToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset password token: %s", err.Error())
	}

	cacheKey := util.ConstructResetPasswordTokenKey(token)
	if err = s.cacheSvc.Set(ctx, cacheKey, user.ID, ResetPasswordTokenTTL); err != nil {
		s.logger.Errorf(
			"userId: %s, failed to set reset password token: %s",
			user.ID, err.Error(),
		)

		return fmt.Errorf("failed to set reset password token: %s", err.Error())
	}

	if err = s.publishUserEvent(ctx, &UserEvent{
		UserID:    user.ID,
		EventType: UserResetPasswordEvent,
		Timestamp: time.Now(),
		Data: map[string]any{
			"email":                user.Email,
			"reset_password_token": token,
		},
	}); err != nil {
		s.logger.Errorf(
			"userId: %s, email: %s, failed to publish user reset password event: %s",
			user.ID, user.Email, err.Error(),
		)

		return fmt.Errorf("failed to publish user reset password event: %s", err.Error())
	}

	return nil
}

func (s *AuthService) ResetPassword(
	ctx context.Context,
	resetPasswordToken int,
	req *reqDto.ResetPasswordRequest,
) error {
	var userID string
	cacheKey := util.ConstructResetPasswordTokenKey(resetPasswordToken)
	if err := s.cacheSvc.Get(ctx, cacheKey, &userID); err != nil {
		if err == redis.Nil {
			return customErr.NewCustomError(customErr.ErrNotFound, "Reset password token not found")
		}

		return fmt.Errorf("failed to get userId from cache: %s", err.Error())
	}

	user, err := s.authRepository.GetById(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customErr.NewCustomError(customErr.ErrNotFound, "User not found")
		}

		return fmt.Errorf("failed to get user by id: %s", err.Error())
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %s", err.Error())
	}

	if err := s.authRepository.UpdateById(ctx, user.ID, &entity.User{
		Password: hashedPassword,
	}); err != nil {
		return fmt.Errorf("failed to update user password: %s", err.Error())
	}

	if err := s.cacheSvc.Delete(ctx, cacheKey); err != nil {
		s.logger.Errorf(
			"userId: %s, failed to delete user reset password cache: %s",
			user.ID, err.Error(),
		)
	}

	return nil
}

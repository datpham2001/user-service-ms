package auth

import (
	"context"
	"time"

	"github.com/datpham/user-service-ms/internal/repository/common"
	"github.com/datpham/user-service-ms/internal/repository/entity"
	"golang.org/x/oauth2"
)

type IAuthRepository interface {
	common.IGenericRepository[entity.User]
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error)
}

type IJwtTokenService interface {
	GenerateTokenPair(userId string) (string, string, error)
}

type IOAuthService interface {
	GetGoogleAuthUrl() string
	VerifyGoogleState(state string) error
	GetGoogleAccessToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetGoogleUserInfo(ctx context.Context, accessToken string) (map[string]any, error)
}

type ICacheService interface {
	Get(ctx context.Context, key string, obj any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

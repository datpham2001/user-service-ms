package auth

import (
	"context"

	"github.com/datpham/user-service-ms/internal/repository/common"
	"github.com/datpham/user-service-ms/internal/repository/entity"
	"golang.org/x/oauth2"
)

type IAuthRepository interface {
	common.IGenericRepository[entity.User]
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type IJwtTokenService interface {
	GenerateTokenPair(userId string) (string, string, error)
}

type IOAuthService interface {
	GetGoogleAuthUrl() string
	VerifyGoogleState(state string) error
	GetGoogleAccessToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetGoogleUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error)
}

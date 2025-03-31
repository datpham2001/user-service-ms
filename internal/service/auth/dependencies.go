package service

import (
	"context"

	"github.com/datpham/user-service-ms/internal/repository/common"
	"github.com/datpham/user-service-ms/internal/repository/entity"
)

type IAuthRepository interface {
	common.IGenericRepository[entity.User]
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type IJwtToken interface {
	GenerateTokenPair(userId string) (string, string, error)
}

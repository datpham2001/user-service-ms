package auth

import (
	"context"

	"github.com/datpham/user-service-ms/internal/repository/common"
	"github.com/datpham/user-service-ms/internal/repository/entity"
	"gorm.io/gorm"
)

type AuthRepository struct {
	*common.GenericRepository[entity.User]
}

func New(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		GenericRepository: common.NewGenericRepository[entity.User](db),
	}
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.GetDB().WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*entity.User, error) {
	var user entity.User
	if err := r.GetDB().WithContext(ctx).
		Where("refresh_token = ?", refreshToken).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

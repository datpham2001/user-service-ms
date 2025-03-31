package common

import (
	"context"

	"gorm.io/gorm"
)

type IGenericRepository[T any] interface {
	GetDB() *gorm.DB
	GetById(ctx context.Context, id string) (*T, error)
	Create(ctx context.Context, entity *T) error
	UpdateById(ctx context.Context, id string, entity *T) error
	DeleteById(ctx context.Context, id string) error
}

package common

import (
	"context"

	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	db *gorm.DB
}

func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db}
}

func (r *GenericRepository[T]) GetDB() *gorm.DB {
	return r.db
}

func (r *GenericRepository[T]) GetById(ctx context.Context, id string) (*T, error) {
	var entity T
	if err := r.GetDB().WithContext(ctx).
		Where("id = ?", id).
		First(&entity).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *GenericRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.GetDB().WithContext(ctx).Create(entity).Error
}

func (r *GenericRepository[T]) UpdateById(ctx context.Context, id string, entity *T) error {
	return r.GetDB().WithContext(ctx).
		Where("id = ?", id).
		Updates(entity).Error
}

func (r *GenericRepository[T]) DeleteById(ctx context.Context, id string) error {
	var entity T
	return r.GetDB().WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity).Error
}

package repository

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
)

type CategoryRepository interface {
	GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Category, *entity.Meta, error)
	FindById(ctx context.Context, id uint) (*entity.Category, error)
	FindByCode(ctx context.Context, code string) (*entity.Category, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, id uint, category *entity.Category) error
	Delete(ctx context.Context, id uint) error
}

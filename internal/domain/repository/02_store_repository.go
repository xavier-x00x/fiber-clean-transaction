package repository

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
)

type StoreRepository interface {
	GetFilter(filter entity.QueryFilter) ([]entity.Store, *entity.Meta, error)
	FindById(id uint) (*entity.Store, error)
	FindByCode(code string) (*entity.Store, error)
	Create(ctx context.Context, store *entity.Store) error
	Update(ctx context.Context, id uint, store *entity.Store) error
	Delete(ctx context.Context, id uint) error
}

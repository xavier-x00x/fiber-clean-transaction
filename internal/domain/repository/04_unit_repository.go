package repository

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/transaction"
)

type UnitRepository interface {
	GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Unit, *entity.Meta, error)
	FindById(ctx context.Context, id uint) (*entity.Unit, error)
	FindByCode(ctx context.Context, code string) (*entity.Unit, error)
	Create(ctx context.Context, tx transaction.Transaction, unit *entity.Unit) error
	Update(ctx context.Context, tx transaction.Transaction, id uint, unit *entity.Unit) error
	Delete(ctx context.Context, tx transaction.Transaction, id uint) error
}

package repository

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
)

type UserRepository interface {
	GetFilter(filter entity.QueryFilter) ([]entity.User, *entity.Meta, error)
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(ID uint) (*entity.User, error)
	Update(ctx context.Context, ID uint, user *entity.User) error
	UpdateAvatar(ID uint, avatar string) error
	Delete(ctx context.Context, ID uint) error
}

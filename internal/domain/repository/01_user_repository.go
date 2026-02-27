package repository

import "fiber-clean-transaction/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindById(id uint) (*entity.User, error)
	UpdateAvatar(id uint, avatar string) error
}

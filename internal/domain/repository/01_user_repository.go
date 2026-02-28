package repository

import "fiber-clean-transaction/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(ID uint) (*entity.User, error)
	UpdateAvatar(ID uint, avatar string) error
}

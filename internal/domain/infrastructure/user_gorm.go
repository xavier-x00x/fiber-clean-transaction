package infrastructure

import (
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"

	"gorm.io/gorm"
)

type UserGormRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserGormRepo{db: db}
}

func (r *UserGormRepo) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *UserGormRepo) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserGormRepo) FindById(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserGormRepo) UpdateAvatar(id uint, avatar string) error {
	return r.db.Model(&entity.User{}).Where("id = ?", id).Update("avatar", avatar).Error
}

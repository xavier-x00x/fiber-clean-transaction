package infrastructure

import (
	"context"
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

func (r *UserGormRepo) GetFilter(filter entity.QueryFilter) ([]entity.User, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.User{})
	return PaginateAndFilter[entity.User](r.db, baseQuery, filter)
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

func (r *UserGormRepo) FindByID(ID uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("id = ?", ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserGormRepo) UpdateAvatar(ID uint, avatar string) error {
	return r.db.Model(&entity.User{}).Where("id = ?", ID).Update("avatar", avatar).Error
}

func (r *UserGormRepo) Update(ctx context.Context, ID uint, data *entity.User) error {
	gormTx := GetDBWithTx(ctx, r.db)
	data.ID = ID
	return gormTx.WithContext(ctx).Updates(data).Error
}

func (r *UserGormRepo) Delete(ctx context.Context, ID uint) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Delete(&entity.User{}, ID).Error
}

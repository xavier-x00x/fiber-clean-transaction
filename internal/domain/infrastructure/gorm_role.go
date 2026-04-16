package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"

	"gorm.io/gorm"
)

type RoleGormRepo struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &RoleGormRepo{db: db}
}

func (r *RoleGormRepo) GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Role, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Role{})
	return PaginateAndFilter[entity.Role](r.db, baseQuery, filter)
}

func (r *RoleGormRepo) FindByID(ctx context.Context, ID uint) (*entity.Role, error) {
	var data entity.Role
	gormTx := GetDBWithTx(ctx, r.db)
	err := gormTx.Preload("Permissions").Where("id = ?", ID).Take(&data).Error
	return &data, err
}

func (r *RoleGormRepo) FindByName(name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("name = ?", name).Take(&role).Error
	return &role, err
}

func (r *RoleGormRepo) Create(ctx context.Context, data *entity.Role) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Create(data).Error
}

func (r *RoleGormRepo) Update(ctx context.Context, ID uint, data *entity.Role) error {
	gormTx := GetDBWithTx(ctx, r.db)
	data.ID = ID
	return gormTx.WithContext(ctx).Updates(data).Error
}

func (r *RoleGormRepo) AssignPermission(ctx context.Context, role *entity.Role, permissions []entity.Permission) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions)
}

func (r *RoleGormRepo) Delete(ctx context.Context, ID uint) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Delete(&entity.Role{}, ID).Error
}

func (r *RoleGormRepo) ClearPermissions(ctx context.Context, role *entity.Role) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(role).Association("Permissions").Clear()
}

func (r *RoleGormRepo) AccessPermission(role_name string, permission_name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Preload("Permissions", "name = ?", permission_name).Where("name = ?", role_name).Take(&role).Error
	return &role, err
}

func (r *RoleGormRepo) GetPermissionsByRole(ctx context.Context, roleName string) ([]entity.Permission, error) {
	var role entity.Role
	err := r.db.Preload("Permissions").Where("name = ?", roleName).Take(&role).Error
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

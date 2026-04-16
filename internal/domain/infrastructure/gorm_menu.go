package infrastructure

import (
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"

	"gorm.io/gorm"
)

type MenuGormRepo struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) repository.MenuRepository {
	return &MenuGormRepo{db: db}
}

func (r *MenuGormRepo) GetAllMenus() ([]entity.NuxtMenu, error) {
	var menus []entity.NuxtMenu
	err := r.db.Find(&menus).Error
	return menus, err
}

func (r *MenuGormRepo) GetMenusByRole(roleName string) ([]entity.NuxtMenu, error) {
	var routes []string
	err := r.db.Table("roles").
		Select("permissions.path").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("roles.name = ?", roleName).
		Pluck("permissions.path", &routes).Error
	if err != nil {
		return nil, err
	}

	var menus []entity.NuxtMenu
	err = r.db.Where("`to` IN ?", routes).Find(&menus).Error
	return menus, err
}

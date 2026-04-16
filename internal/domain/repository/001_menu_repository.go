package repository

import (
	"fiber-clean-transaction/internal/domain/entity"
)

type MenuRepository interface {
	GetAllMenus() ([]entity.NuxtMenu, error)
	GetMenusByRole(roleName string) ([]entity.NuxtMenu, error)
}

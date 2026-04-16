package usecase

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"
	"fiber-clean-transaction/pkg/utils"
)

type MenuUsecase struct {
	repo repository.MenuRepository
}

func NewMenuUsecase(repo repository.MenuRepository) *MenuUsecase {
	return &MenuUsecase{
		repo: repo,
	}
}

func (u *MenuUsecase) GetMenusByRole(ctx context.Context, roleName string) ([]entity.NuxtMenu, error) {

	allowedRoles := []string{"admin", "superadmin", "programmer"}
	if utils.Contains(allowedRoles, roleName) {
		menus, err := u.repo.GetAllMenus()
		if err != nil {
			return nil, utils.Internal("Failed to get all menus: "+err.Error(), err)
		}
		return menus, nil
	}

	menus, err := u.repo.GetMenusByRole(roleName)
	if err != nil {
		return nil, utils.Internal("Failed to get menus: "+err.Error(), err)
	}
	return menus, nil
}

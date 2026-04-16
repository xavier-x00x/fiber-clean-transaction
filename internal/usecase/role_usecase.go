package usecase

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	repo "fiber-clean-transaction/internal/domain/repository"
	"fiber-clean-transaction/internal/dto"
	trans "fiber-clean-transaction/internal/transaction"
	"fiber-clean-transaction/pkg/utils"
	"fiber-clean-transaction/pkg/validation"
	"log"
)

type RoleUsecase struct {
	validator *validation.ValidatorHelper
	uow       trans.UnitOfWork
	roleRepo  repo.RoleRepository
	permRepo  repo.PermissionRepository
}

type RoleUsecaseDeps struct {
	Validator *validation.ValidatorHelper
	UOW       trans.UnitOfWork
	RoleRepo  repo.RoleRepository
	PermRepo  repo.PermissionRepository
	// Huruf depan diawali dengan huruf besar
	// Agar bisa diakses di luar package
}

func NewRoleUsecase(deps RoleUsecaseDeps) *RoleUsecase {
	// dari registration di container dimasukan ke struct
	return &RoleUsecase{
		validator: deps.Validator,
		uow:       deps.UOW,
		roleRepo:  deps.RoleRepo,
		permRepo:  deps.PermRepo, // permission repository
	}
}

func (u *RoleUsecase) GetAllFilter(ctx context.Context, meta *dto.MetaRequest) ([]entity.Role, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "updated_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)

	data, resMeta, err := u.roleRepo.GetAllFilter(ctx, filter)
	if err != nil {
		return nil, nil, utils.Internal(err.Error(), err)
	}
	return data, resMeta, nil
}

func (u *RoleUsecase) FindByID(ctx context.Context, ID uint) (*entity.Role, error) {
	data, err := u.roleRepo.FindByID(ctx, ID)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	log.Printf("Role found: %+v", data)
	return data, nil
}

func (u *RoleUsecase) FindByName(name string) (*entity.Role, error) {
	data, err := u.roleRepo.FindByName(name)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *RoleUsecase) Create(ctx context.Context, request *dto.RoleRequest) error {
	// Validasi SEBELUM transaksi
	if err := u.validator.Validate(request); err != nil {
		return err
	}

	return u.uow.Do(ctx, func(ctx context.Context) error {
		role := &entity.Role{
			Name: request.Name,
		}

		// validasi permissions + mapping
		permissions := make([]entity.Permission, 0)
		for _, permCode := range request.Permissions {
			perm, err := u.permRepo.FindByName(permCode)
			if err != nil {
				return utils.NotFound("permission " + permCode + " not found")
			}
			permissions = append(permissions, *perm)
		}

		if err := u.roleRepo.Create(ctx, role); err != nil {
			return utils.Internal(err.Error(), err)
		}

		if err := u.roleRepo.AssignPermission(ctx, role, permissions); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil
	})
}

func (u *RoleUsecase) Update(ctx context.Context, ID uint, request *dto.RoleRequest) error {
	// Validasi SEBELUM transaksi
	if err := u.validator.ValidateUpdate(request, ID); err != nil {
		return err
	}

	return u.uow.Do(ctx, func(ctx context.Context) error {
		if _, err := u.roleRepo.FindByID(ctx, ID); err != nil {
			return utils.NotFound(err.Error())
		}

		// validasi permissions + mapping
		permissions := make([]entity.Permission, 0)
		for _, permCode := range request.Permissions {
			perm, err := u.permRepo.FindByName(permCode)
			if err != nil {
				return utils.NotFound("permission " + permCode + " not found")
			}
			permissions = append(permissions, *perm)
		}

		role := &entity.Role{
			Name: request.Name,
		}

		if err := u.roleRepo.Update(ctx, ID, role); err != nil {
			return utils.Internal(err.Error(), err)
		}

		if err := u.roleRepo.AssignPermission(ctx, role, permissions); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil
	})
}

func (u *RoleUsecase) Delete(ctx context.Context, ID uint) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		if _, err := u.roleRepo.FindByID(ctx, ID); err != nil {
			return utils.NotFound(err.Error())
		}

		if err := u.roleRepo.Delete(ctx, ID); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil
	})
}

func (u *RoleUsecase) Authorization(role_name string, permission_name string) (bool, error) {

	allowedRoles := []string{"admin", "superadmin", "programmer"}
	if utils.Contains(allowedRoles, role_name) {
		return true, nil
	}

	data, err := u.roleRepo.AccessPermission(role_name, permission_name)

	if err != nil {
		return false, utils.Internal(err.Error(), err)
	}

	if data.ID == 0 {
		return false, nil
	}

	return true, nil
}

func (u *RoleUsecase) GetPermissionsByRole(ctx context.Context, roleName string) ([]entity.Permission, error) {

	allowedRoles := []string{"admin", "superadmin", "programmer"}

	if utils.Contains(allowedRoles, roleName) {
		// Jika role termasuk dalam allowedRoles, kembalikan semua permissions
		permissions, err := u.permRepo.GetAllPermissions(ctx)
		if err != nil {
			return nil, utils.Internal(err.Error(), err)
		}
		return permissions, nil
	}

	permissions, err := u.roleRepo.GetPermissionsByRole(ctx, roleName)
	if err != nil {
		return nil, utils.Internal(err.Error(), err)
	}
	return permissions, nil
}

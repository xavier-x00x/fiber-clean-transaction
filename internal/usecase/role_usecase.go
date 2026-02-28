package usecase

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/transaction"
	"fiber-clean-transaction/pkg/utils"
	"fiber-clean-transaction/pkg/validation"
)

type RoleUsecase struct {
	validator *validation.ValidatorHelper
	uow       transaction.UnitOfWork
	repo      repository.RoleRepository
}

func NewRoleUsecase(r repository.RoleRepository, uow transaction.UnitOfWork, v *validation.ValidatorHelper) *RoleUsecase {
	return &RoleUsecase{
		validator: v,
		uow:       uow,
		repo:      r,
	}
}

func (u *RoleUsecase) GetAllFilter(ctx context.Context, meta *dto.MetaRequest) ([]entity.Role, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "updated_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)

	data, resMeta, err := u.repo.GetAllFilter(ctx, filter)
	if err != nil {
		return nil, nil, utils.Internal(err.Error(), err)
	}
	return data, resMeta, nil
}

func (u *RoleUsecase) FindByID(ctx context.Context, ID uint) (*entity.Role, error) {
	data, err := u.repo.FindByID(ctx, ID)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *RoleUsecase) FindByName(name string) (*entity.Role, error) {
	data, err := u.repo.FindByName(name)
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

		if err := u.repo.Create(ctx, role); err != nil {
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
		if _, err := u.repo.FindByID(ctx, ID); err != nil {
			return utils.NotFound(err.Error())
		}

		role := &entity.Role{
			Name: request.Name,
		}

		if err := u.repo.Update(ctx, ID, role); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil
	})
}

func (u *RoleUsecase) Delete(ctx context.Context, ID uint) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		if _, err := u.repo.FindByID(ctx, ID); err != nil {
			return utils.NotFound(err.Error())
		}

		if err := u.repo.Delete(ctx, ID); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil
	})
}

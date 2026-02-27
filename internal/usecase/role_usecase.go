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
	repo      repository.RoleRepository
	txManager transaction.UseCaseRunner
	validator *validation.ValidatorHelper
}

func NewRoleUsecase(r repository.RoleRepository, t transaction.UseCaseRunner, v *validation.ValidatorHelper) *RoleUsecase {
	return &RoleUsecase{
		repo:      r,
		txManager: t,
		validator: v,
	}
}

func (u *RoleUsecase) GetAllFilter(meta *dto.MetaRequest) ([]entity.Role, *entity.Meta, error) {

	// validasi order by untuk hindari SQL injection
	direction := []string{"asc", "desc"}
	order := []string{"id", "code", "name", "updated_at"}
	searchColumn := []string{"id", "code", "name"}

	if !utils.Contains(order, meta.OrderColumn) || !utils.Contains(direction, meta.OrderDir) {
		meta.OrderColumn = "id"
		meta.OrderDir = "asc"
	}

	filter := entity.QueryFilter{
		Page:         meta.Page,
		Limit:        meta.Limit,
		Search:       meta.Search,
		OrderColumn:  meta.OrderColumn,
		OrderDir:     meta.OrderDir,
		SearchColumn: searchColumn,
		Conditions:   map[string]interface{}{},
	}

	filter.Conditions["status"] = 1

	data, resMeta, err := u.repo.GetFilter(filter)
	if err != nil {
		return nil, nil, utils.Internal(err.Error(), err)
	}
	return data, resMeta, nil
}

func (u *RoleUsecase) FindById(id uint) (*entity.Role, error) {
	data, err := u.repo.FindById(id)
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

	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.Validate(request); err != nil {
		return err
	}

	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	txErr := u.txManager.WithTransaction(ctx, func(ctx context.Context) error {

		Role := &entity.Role{
			Name: request.Name,
		}

		if err := u.repo.Create(ctx, Role); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})

	return txErr
}

func (u *RoleUsecase) Update(ctx context.Context, id uint, request *dto.RoleRequest) error {

	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.ValidateUpdate(request, id); err != nil {
		return err
	}

	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	txErr := u.txManager.WithTransaction(ctx, func(ctx context.Context) error {

		// cek apakah data ada
		if _, err := u.repo.FindById(id); err != nil {
			return utils.NotFound(err.Error())
		}

		Role := &entity.Role{
			Name: request.Name,
		}

		if err := u.repo.Update(ctx, id, Role); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})

	return txErr
}

func (u *RoleUsecase) Delete(ctx context.Context, id uint) error {
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	txErr := u.txManager.WithTransaction(ctx, func(ctx context.Context) error {

		// cek apakah data ada
		if _, err := u.repo.FindById(id); err != nil {
			return utils.NotFound(err.Error())
		}

		// hapus data
		if err := u.repo.Delete(ctx, id); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil // berhasil di commit
	})

	return txErr
}

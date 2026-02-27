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

type PermissionUsecase struct {
	validator *validation.ValidatorHelper
	uow       transaction.UnitOfWork
	repo      repository.PermissionRepository
}

func NewPermissionUsecase(r repository.PermissionRepository, uow transaction.UnitOfWork, v *validation.ValidatorHelper) *PermissionUsecase {
	return &PermissionUsecase{
		validator: v,
		uow:       uow,
		repo:      r,
	}
}

func (u *PermissionUsecase) GetAllFilter(ctx context.Context, meta *dto.MetaRequest) ([]entity.Permission, *entity.Meta, error) {
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

	data, resMeta, err := u.repo.GetAllFilter(ctx, filter)
	if err != nil {
		return nil, nil, utils.Internal(err.Error(), err)
	}
	return data, resMeta, nil
}

func (u *PermissionUsecase) FindById(ctx context.Context, id uint) (*entity.Permission, error) {
	data, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *PermissionUsecase) Create(ctx context.Context, request *dto.PermissionRequest) error {
	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.Validate(request); err != nil {
		return err
	}
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {
		data := &entity.Permission{
			Path: request.Path,
			Name: request.Name,
		}

		if err := u.repo.Create(ctx, data); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})
}

func (u *PermissionUsecase) Update(ctx context.Context, id uint, request *dto.PermissionRequest) error {
	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.ValidateUpdate(request, id); err != nil {
		return err
	}
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {

		// cek apakah data ada
		if _, err := u.repo.FindById(ctx, id); err != nil {
			return utils.NotFound(err.Error())
		}

		data := &entity.Permission{
			Path: request.Path,
			Name: request.Name,
		}

		if err := u.repo.Update(ctx, id, data); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})
}

func (u *PermissionUsecase) SyncPermissions(ctx context.Context, permissions []*dto.PermissionRequest) error {
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {

		for _, perm := range permissions {
			existingPerm, err := u.repo.FindByCode(ctx, perm.Path)
			if err != nil {
				return utils.Internal(err.Error(), err)
			}
			if existingPerm == nil {
				newPerm := &entity.Permission{
					Path: perm.Path,
					Name: perm.Name,
				}
				if err := u.repo.Create(ctx, newPerm); err != nil {
					return utils.Internal(err.Error(), err)
				}
			}
		}

		return nil // berhasil di commit
	})
}

func (u *PermissionUsecase) Delete(ctx context.Context, id uint) error {
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {

		// cek apakah data ada
		if _, err := u.repo.FindById(ctx, id); err != nil {
			return utils.NotFound(err.Error())
		}

		// hapus data
		if err := u.repo.Delete(ctx, id); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil // berhasil di commit
	})
}

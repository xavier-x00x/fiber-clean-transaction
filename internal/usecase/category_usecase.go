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

type CategoryUsecase struct {
	repo      repository.CategoryRepository
	txManager transaction.UseCaseRunner
	validator *validation.ValidatorHelper
}

func NewCategoryUsecase(r repository.CategoryRepository, t transaction.UseCaseRunner, v *validation.ValidatorHelper) *CategoryUsecase {
	return &CategoryUsecase{
		repo:      r,
		txManager: t,
		validator: v,
	}
}

func (u *CategoryUsecase) GetAllFilter(ctx context.Context, meta *dto.MetaRequest) ([]entity.Category, *entity.Meta, error) {
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

	data, resMeta, err := u.repo.GetAllFilter(ctx, filter)
	if err != nil {
		return nil, nil, utils.Internal(err.Error(), err)
	}
	return data, resMeta, nil
}

func (u *CategoryUsecase) FindById(ctx context.Context, id uint) (*entity.Category, error) {
	data, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *CategoryUsecase) FindByCode(ctx context.Context, code string) (*entity.Category, error) {
	data, err := u.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *CategoryUsecase) Create(ctx context.Context, request *dto.CategoryRequest) error {
	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.Validate(request); err != nil {
		return err
	}
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	txErr := u.txManager.WithTransaction(ctx, func(ctx context.Context) error {

		category := &entity.Category{
			StoreCode:   request.StoreCode,
			Code:        request.Code,
			Name:        request.Name,
			Description: request.Description,
			Status:      request.Status,
		}

		if err := u.repo.Create(ctx, category); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})

	return txErr
}

func (u *CategoryUsecase) Update(ctx context.Context, id uint, request *dto.CategoryRequest) error {
	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.ValidateUpdate(request, id); err != nil {
		return err
	}
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	txErr := u.txManager.WithTransaction(ctx, func(ctx context.Context) error {

		// cek apakah data ada
		if _, err := u.repo.FindById(ctx, id); err != nil {
			return utils.NotFound(err.Error())
		}

		category := &entity.Category{
			StoreCode:   request.StoreCode,
			Code:        request.Code,
			Name:        request.Name,
			Description: request.Description,
			Status:      request.Status,
		}

		if err := u.repo.Update(ctx, id, category); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})

	return txErr
}

func (u *CategoryUsecase) Delete(ctx context.Context, id uint) error {
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	txErr := u.txManager.WithTransaction(ctx, func(ctx context.Context) error {

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

	return txErr
}

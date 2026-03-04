package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NumberSequenceGormRepo struct {
	db *gorm.DB
}

func NewNumberSequenceRepository(db *gorm.DB) repository.NumberSequenceRepository {
	return &NumberSequenceGormRepo{db: db}
}

func (r *NumberSequenceGormRepo) GetNextNumber(ctx context.Context, prefix string, period string) (int, error) {
	gormTx := GetDBWithTx(ctx, r.db)

	var seq entity.NumberSequence

	// Gunakan row-level locking (SELECT ... FOR UPDATE) agar aman concurrent
	err := gormTx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("prefix = ? AND period = ?", prefix, period).
		First(&seq).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Belum ada record, buat baru dengan last_number = 1
			seq = entity.NumberSequence{
				Prefix:     prefix,
				Period:     period,
				LastNumber: 1,
			}
			if err := gormTx.WithContext(ctx).Create(&seq).Error; err != nil {
				return 0, err
			}
			return seq.LastNumber, nil
		}
		return 0, err
	}

	// Increment last_number
	seq.LastNumber++
	if err := gormTx.WithContext(ctx).
		Model(&seq).
		Update("last_number", seq.LastNumber).Error; err != nil {
		return 0, err
	}

	return seq.LastNumber, nil
}

package infrastructure

import (
	"context"
	domainTx "fiber-clean-transaction/internal/transaction"

	"gorm.io/gorm"
)

type GormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) domainTx.TransactionManager {
	return &GormTransactionManager{db: db}
}

func (tm *GormTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx domainTx.Transaction) error) error {
	return tm.db.WithContext(ctx).Transaction(func(gormTx *gorm.DB) error {
		// Konversi *gorm.DB ke abstraksi Transaction
		return fn(ctx, gormTx)
	})
}

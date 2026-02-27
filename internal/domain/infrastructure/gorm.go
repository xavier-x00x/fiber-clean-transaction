package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/transaction"

	"gorm.io/gorm"
)

func GetDBWithTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := transaction.GetTx(ctx); ok {
		return tx
	}
	return db
}

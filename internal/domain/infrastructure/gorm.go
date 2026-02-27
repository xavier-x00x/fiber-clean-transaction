package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/transaction"

	"gorm.io/gorm"
)

func GetDBWithTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	return transaction.DBFromContext(ctx, db)
}

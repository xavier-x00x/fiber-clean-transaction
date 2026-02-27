package transaction

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

var transactionKey = txKey{}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

func GetTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(transactionKey).(*gorm.DB)
	return tx, ok
}

type GormRunner struct {
	DB *gorm.DB
}

func NewGormRunner(db *gorm.DB) UseCaseRunner {
	return &GormRunner{DB: db}
}

func (r *GormRunner) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// inject tx ke context
		ctx = WithTx(ctx, tx)

		return fn(ctx)
	})
}

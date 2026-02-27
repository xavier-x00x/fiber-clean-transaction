package transaction

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

type trxKey struct{}

type GormUnitOfWork struct {
	db         *gorm.DB
	maxRetry   int
	retryDelay time.Duration
}

func NewGormUnitOfWork(db *gorm.DB) *GormUnitOfWork {
	return &GormUnitOfWork{
		db:         db,
		maxRetry:   3,
		retryDelay: 100 * time.Millisecond, // 100ms
	}
}

func (u *GormUnitOfWork) Do(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {

	var err error

	for attempt := 0; attempt <= u.maxRetry; attempt++ {

		err = u.db.WithContext(ctx).
			Transaction(func(tx *gorm.DB) error {

				ctxWithTx :=
					context.WithValue(ctx, trxKey{}, tx)

				return fn(ctxWithTx)
			})

		if err == nil {
			return nil
		}

		// bukan deadlock → langsung gagal
		if !isDeadlockError(err) {
			return err
		}

		// retry delay (exponential backoff)
		sleep := u.retryDelay * time.Duration(attempt+1)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleep):
		}
	}

	return err
}

func isDeadlockError(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	return strings.Contains(msg, "Error 1213") || // deadlock
		strings.Contains(msg, "Error 1205") || // lock timeout
		strings.Contains(msg, "deadlock")
}

func DBFromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(trxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return db
}

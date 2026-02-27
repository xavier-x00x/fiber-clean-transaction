package transaction

import "context"

// Transaction adalah abstraksi yang tidak bergantung pada implementasi database
type Transaction interface{}

// TransactionManager adalah interface untuk mengelola transaction
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}

type UseCaseRunner interface {
	WithTransaction(
		ctx context.Context,
		fn func(ctx context.Context) error,
	) error
}

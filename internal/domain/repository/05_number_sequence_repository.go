package repository

import "context"

type NumberSequenceRepository interface {
	// GetNextNumber retrieves and increments the next sequence number
	// for the given prefix and period (YYMM).
	// If no record exists yet, it creates one starting at 1.
	// Uses row-level locking for concurrency safety.
	GetNextNumber(ctx context.Context, prefix string, period string) (int, error)
}

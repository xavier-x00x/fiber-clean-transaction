package docnumber

import (
	"context"
	"fiber-clean-transaction/internal/domain/repository"
	"fmt"
	"time"
)

// GenerateDocNumber generates a document/invoice number with the format:
// <PREFIX><YYMM>-<5_DIGIT_NUMBER>
//
// Example output: "INV2603-00001", "PO2603-00012"
//
// Parameters:
//   - ctx: context (should contain transaction if called within UoW)
//   - repo: NumberSequenceRepository for retrieving/incrementing sequence
//   - prefix: document prefix, e.g. "INV", "PO", "SO"
//
// The period (YYMM) is derived from the current time automatically.
func GenerateDocNumber(ctx context.Context, repo repository.NumberSequenceRepository, prefix string) (string, error) {
	period := time.Now().Format("0601") // YYMM format, e.g. "2603" for March 2026

	nextNum, err := repo.GetNextNumber(ctx, prefix, period)
	if err != nil {
		return "", fmt.Errorf("failed to generate doc number: %w", err)
	}

	// Format: <PREFIX><YYMM>-<5_DIGIT>
	docNumber := fmt.Sprintf("%s%s-%05d", prefix, period, nextNum)
	return docNumber, nil
}

package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/miztch/sasha/internal/domain"
)

// MatchWriter is an interface for writing matches
type StdoutMatchWriter struct{}

// WriteMatches writes matches to stdout
func (w *StdoutMatchWriter) WriteMatches(ctx context.Context, matches []domain.Match) error {
	if matches == nil {
		matches = make([]domain.Match, 0)
	}

	m, err := json.Marshal(matches)
	if err != nil {
		return fmt.Errorf("failed to marshal matches: %w", err)
	}

	fmt.Println(string(m))

	return nil
}

package infrastructure

import (
	"context"

	"github.com/miztch/sasha/internal/domain"
)

// MatchWriter is an interface for writing matches
type MatchWriter interface {
	WriteMatches(ctx context.Context, matches []domain.Match) error
}

package infrastructure

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/miztch/sasha/internal/domain"
)

// MatchWriter is an interface for writing matches
type matchRepositoryImpl struct {
	vlrGGScraper *VlrGGScraper
	matchWriter  MatchWriter
}

// NewMatchRepository creates a new match repository
func NewMatchRepository(vlrGGScraper *VlrGGScraper, matchWriter MatchWriter) *matchRepositoryImpl {
	return &matchRepositoryImpl{
		vlrGGScraper: vlrGGScraper,
		matchWriter:  matchWriter,
	}
}

// ScrapeMatch scrapes a match
func (r *matchRepositoryImpl) ScrapeMatch(matchUrlPath string) (domain.VlrMatch, error) {
	match, err := r.vlrGGScraper.scrapeMatch(matchUrlPath)
	if err != nil {
		return domain.VlrMatch{}, fmt.Errorf("failed to scrape match: %w", err)
	}
	slog.Info(fmt.Sprintf("Scraped match: %d", match.Id), "Error", "")
	return match, nil
}

// GetMatchURLList gets a list of match URLs
func (r *matchRepositoryImpl) GetMatchURLList(pageNumber int) ([]string, error) {
	matchURLs, err := r.vlrGGScraper.getMatchURLList(pageNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get match urls: %w", err)
	}

	return matchURLs, nil
}

// WriteMatches writes matches
func (r *matchRepositoryImpl) WriteMatches(ctx context.Context, matches []domain.Match) error {
	err := r.matchWriter.WriteMatches(ctx, matches)

	if err != nil {
		return fmt.Errorf("failed to write matches: %w", err)
	}

	return nil
}

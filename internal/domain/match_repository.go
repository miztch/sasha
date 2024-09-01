package domain

import "context"

// MatchRepository is an interface for match repositories
type MatchRepository interface {
	// GetMatchURLList gets a list of match URLs
	GetMatchURLList(page int) ([]string, error)
	// ScrapeMatch scrapes a match
	ScrapeMatch(matchURL string) (VlrMatch, error)
	// WriteMatches writes matches
	WriteMatches(ctx context.Context, matches []Match) error
}

// if the VlrMatch is empty, return true
func IsEmptyVlrMatch(m VlrMatch) bool {
	if m.Id != 0 && m.PagePath != "" && m.EventPagePath != "" {
		return false
	}

	if m.StartTime != "" || m.StartDate != "" {
		return false
	}

	if m.Name != "" || m.BestOf != 0 || len(m.Teams) > 0 {
		return false
	}

	return true
}

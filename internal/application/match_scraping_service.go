package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/miztch/sasha/internal/domain"
)

type MatchService struct {
	matchRepo domain.MatchRepository
	eventRepo domain.EventRepository
}

func NewMatchService(matchRepo domain.MatchRepository, eventRepo domain.EventRepository) *MatchService {
	return &MatchService{
		matchRepo: matchRepo,
		eventRepo: eventRepo,
	}
}

func (svc *MatchService) FetchMatches(page int) ([]domain.Match, error) {
	// マッチの取得ロジックを実装
	matchURLs, err := svc.matchRepo.GetMatchURLList(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get match urls: %w", err)
	}

	var matches []domain.Match
	for _, matchURL := range matchURLs {
		m, err := svc.matchRepo.ScrapeMatch(matchURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get match: %w", err)
		}

		// Check if the scraped match is empty
		if domain.IsEmptyVlrMatch(m) {
			slog.Warn("skipping empty match", "matchURL", matchURL)
			continue
		}

		e, err := svc.eventRepo.GetEvent(m.EventPagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get event: %w", err)
		}

		match := domain.NewMatch(m, e)
		matches = append(matches, match)
	}

	return matches, nil
}

func (svc *MatchService) WriteMatches(ctx context.Context, matches []domain.Match) error {
	// マッチの書き込みロジックを実装
	err := svc.matchRepo.WriteMatches(ctx, matches)
	if err != nil {
		return fmt.Errorf("failed to write matches: %w", err)
	}

	return nil
}

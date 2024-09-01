package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/miztch/sasha/internal/application"
	"github.com/miztch/sasha/internal/config"
	"github.com/miztch/sasha/internal/infrastructure"
)

// Payload is the input data structure
type Payload struct {
	Page int `json:"page"`
}

// Response is the output data structure
type Response struct {
	MatchesCount int `json:"matches_count"`
}

// IsRunningOnLambda checks if the code is running on AWS Lambda
func IsRunningOnLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

// NewApp creates a new application
func NewApp() (*application.MatchService, error) {
	matchWriter, err := config.NewMatchWriter()
	if err != nil {
		return nil, fmt.Errorf("failed to create match writer: %v", err)
	}

	service := application.NewMatchService(
		infrastructure.NewMatchRepository(infrastructure.NewVlrGGScraper(), matchWriter),
		infrastructure.NewEventRepository(infrastructure.NewVlrGGScraper(), infrastructure.NewEventCache()),
	)
	return service, nil
}

// Run is the main function
func Run(payload Payload) (string, error) {
	app, err := NewApp()
	if err != nil {
		slog.Error("failed to create app", "Error", err.Error())
		return "", err
	}

	matches, err := app.FetchMatches(payload.Page)
	if err != nil {
		slog.Error("failed to fetch matches", "Error", err.Error())
		return "", err
	}

	err = app.WriteMatches(context.Background(), matches)
	if err != nil {
		slog.Error("failed to write matches", "Error", err.Error())
		return "", err
	}

	r, _ := json.Marshal(Response{MatchesCount: len(matches)})

	return string(r), nil
}

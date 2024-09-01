package config

import (
	"context"
	"fmt"
	"os"

	"github.com/miztch/sasha/internal/infrastructure"
)

// NewMatchWriter creates a new match writer
func NewMatchWriter() (infrastructure.MatchWriter, error) {
	switch outputMode := os.Getenv("OUTPUT_MODE"); outputMode {
	case "", "stdout":
		return &infrastructure.StdoutMatchWriter{}, nil
	case "dynamodb":
		dbClient, err := infrastructure.NewDynamoDBClient(context.Background(), GetDynamoDBConfig().TableName)
		if err != nil {
			return nil, fmt.Errorf("failed to create DynamoDB client: %v", err)
		}
		return dbClient, nil
	default:
		return nil, fmt.Errorf("invalid output mode: %s", outputMode)
	}
}

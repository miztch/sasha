package infrastructure

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/miztch/sasha/internal/domain"
)

// DynamoDB BatchWriteItem has a limit of 25 items
const batchSize = 25

// DynamoDBClient is a client for DynamoDB
type DynamoDBClient struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBClient creates a new DynamoDBClient
func NewDynamoDBClient(ctx context.Context, tableName string) (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %w", err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBClient{client: client, tableName: tableName}, nil
}

// BatchWriteMatch writes matches to DynamoDB
func (d *DynamoDBClient) BatchWriteMatch(ctx context.Context, writeReqs []types.WriteRequest) error {
	_, err := d.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{d.tableName: writeReqs},
	})
	if err != nil {
		return fmt.Errorf("failed to write batch: %w", err)
	}
	return nil
}

// MatchBatchWriter writes matches in batches
// This is a workaround for the 25-item limit of BatchWriteItem
func (d *DynamoDBClient) MatchBatchWriter(ctx context.Context, matches []domain.Match) error {
	var allWriteReqs []types.WriteRequest

	// Prepare all write requests
	for _, match := range matches {
		item, err := attributevalue.MarshalMap(match)
		if err != nil {
			return fmt.Errorf("failed to marshal match: %w", err)
		}
		allWriteReqs = append(allWriteReqs, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	}

	// Write in batches
	for i := 0; i < len(allWriteReqs); i += batchSize {
		end := i + batchSize
		if end > len(allWriteReqs) {
			end = len(allWriteReqs)
		}

		chunk := allWriteReqs[i:end]

		err := d.BatchWriteMatch(ctx, chunk)
		if err != nil {
			return fmt.Errorf("failed to write batch: %w", err)
		}
	}
	return nil
}

// for matchWriter interface
func (d *DynamoDBClient) WriteMatches(ctx context.Context, matches []domain.Match) error {
	return d.MatchBatchWriter(ctx, matches)
}

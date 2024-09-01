package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvFromConfig loads environment variables from a .env file
func LoadEnvFromConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	return nil
}

// dynamoDBConfig is a configuration for DynamoDB
type dynamoDBConfig struct {
	TableName   string
	Region      *string
	EndpointURL *string
}

// getDynamoDBConfig gets a DynamoDB configuration from environment variables
func GetDynamoDBConfig() dynamoDBConfig {
	region := os.Getenv("AWS_REGION")
	endpointURL := os.Getenv("DYNAMODB_ENDPOINT_URL")
	return dynamoDBConfig{
		TableName:   os.Getenv("VLR_MATCHES_TABLE"),
		Region:      &region,
		EndpointURL: &endpointURL,
	}
}

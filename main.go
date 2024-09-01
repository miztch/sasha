package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/miztch/sasha/internal"
	"github.com/miztch/sasha/internal/config"
)

func Handler(ctx context.Context, payload internal.Payload) (string, error) {
	response, _ := internal.Run(payload)

	return response, nil
}

func init() {
	// Load environment variables from config
	config.LoadEnvFromConfig()

	// Set up logging
	if internal.IsRunningOnLambda() {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	} else {
		logFile, err := os.OpenFile("sasha.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("failed to open log file", "Error", err)
		}
		slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, nil)))
	}
}

func main() {
	if internal.IsRunningOnLambda() {
		lambda.Start(Handler)
	} else {
		pageNum := flag.Int("page", 1, "Page number")
		flag.Parse()
		Handler(context.TODO(), internal.Payload{Page: *pageNum})
	}
}

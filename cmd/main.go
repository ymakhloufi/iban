package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ymakhloufi/pfc/internal/http"
	"github.com/ymakhloufi/pfc/internal/pkg/iban"
	"go.uber.org/zap"
)

const port = 80

func main() {
	logger, err := newLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	ibanService := iban.NewService()
	ibanController := iban.NewController(ibanService, logger)
	httpServer := http.NewHttpServer(port, logger, []http.Controller{
		ibanController,
	})

	logger.Info(fmt.Sprintf("Starting Web Server on port %d", port))
	err = httpServer.Run()
	logger.Warn("http server stopped", zap.Error(err))
}

// newLogger creates a new logger, depending on the environment variable ENVIRONMENT.
func newLogger() (*zap.Logger, error) {
	if os.Getenv("ENVIRONMENT") == "dev" {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("failed to create zap development logger: %w", err)
		}
		return logger, nil
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create zap production logger: %w", err)
	}
	return logger, nil
}

// Filename: cmd/api/main.go
package main

import (
	"flag"
	"log/slog"
	"os"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port        int
	environment string
}

type applicationDependencies struct {
	config serverConfig
	logger *slog.Logger
}

func main() {
	var settings serverConfig
	
	flag.IntVar(&settings.port, "port", 4000, "Server port")
	flag.StringVar(&settings.environment, "env", "development",
		"Environment(development|staging|production)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	appInstance := &applicationDependencies{
		config: settings,
		logger: logger,
	}

	// Use the serve method which properly handles routing
	err := appInstance.serve()
	if err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
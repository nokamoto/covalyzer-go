package main

import (
	"log/slog"
	"os"

	"github.com/nokamoto/covalyzer-go/internal/infra/config"
	"github.com/nokamoto/covalyzer-go/internal/infra/gh"
	"github.com/nokamoto/covalyzer-go/internal/usecase"
)

func main() {
	config, err := config.NewConfig("config.yaml")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	gh := gh.NewGitHub()
	err = usecase.NewCovalyzer(config, gh).Run()
	if err != nil {
		slog.Error("failed to run", "error", err)
		os.Exit(1)
	}
}

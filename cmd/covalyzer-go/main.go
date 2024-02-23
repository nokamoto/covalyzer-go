package main

import (
	"log/slog"
	"os"

	"github.com/nokamoto/covalyzer-go/internal/infra/config"
	"github.com/nokamoto/covalyzer-go/internal/infra/gh"
	"github.com/nokamoto/covalyzer-go/internal/usecase"
)

func main() {
	level := slog.LevelInfo
	if os.Getenv("DEBUG") != "" {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))

	config, err := config.NewConfig("config.yaml")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	gh, err := gh.NewGitHub()
	if err != nil {
		slog.Error("failed to create a GitHub client", "error", err)
		os.Exit(1)
	}
	err = usecase.NewCovalyzer(config, gh).Run()
	if err != nil {
		slog.Error("failed to run", "error", err)
		os.Exit(1)
	}
}

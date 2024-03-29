package main

import (
	"log/slog"
	"os"

	"github.com/nokamoto/covalyzer-go/internal/infra/command"
	"github.com/nokamoto/covalyzer-go/internal/infra/config"
	"github.com/nokamoto/covalyzer-go/internal/infra/writer"
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

	configYAML := os.Getenv("CONFIG_YAML")
	if configYAML == "" {
		configYAML = "config.yaml"
	}

	config, err := config.NewConfig(configYAML)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	wd, err := command.NewWorkingDir()
	if err != nil {
		slog.Error("failed to create a working directory", "error", err)
		os.Exit(1)
	}

	gh, err := command.NewGitHub(wd)
	if err != nil {
		slog.Error("failed to create a GitHub client", "error", err)
		os.Exit(1)
	}
	gt := command.NewGoTool(wd)

	res, err := usecase.NewCovalyzer(config, gh, gt).Run()
	if err != nil {
		slog.Error("failed to run", "error", err)
		os.Exit(1)
	}
	slog.Info("coverage", "coverage", res)
	wd.Clean()

	w := writer.NewCSVWriter()
	if err := w.Write(config, res); err != nil {
		slog.Error("failed to write", "error", err)
		os.Exit(1)
	}
}

//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock_test.go -package=$GOPACKAGE
package usecase

import (
	"fmt"
	"log/slog"

	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

type gh interface {
	Clone(*v1.Repository) error
}

type Covalyzer struct {
	config *v1.Config
	gh     gh
}

func NewCovalyzer(config *v1.Config, gh gh) *Covalyzer {
	return &Covalyzer{
		config: config,
		gh:     gh,
	}
}

func (c *Covalyzer) Run() error {
	for _, repo := range c.config.GetRepositories() {
		logger := slog.With("repo", repo)
		logger.Info("run")
		if err := c.gh.Clone(repo); err != nil {
			logger.Error("failed to clone", "error", err)
			return fmt.Errorf("failed to clone %v: %w", repo, err)
		}
	}
	return nil
}

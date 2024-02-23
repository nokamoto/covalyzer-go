//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock_test.go -package=$GOPACKAGE
package usecase

import (
	"errors"
	"fmt"
	"log/slog"

	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

var (
	ErrCommitNotFound = fmt.Errorf("commit not found")
)

type gh interface {
	// Clone clones a repository and returns the path to the cloned repository.
	Clone(*v1.Repository) (string, error)
	// Checkout checks out a repository at a specific timestamp and returns the commit.
	// The commit is the most recent commit before the timestamp.
	//
	// If the commit is not found for the timestamp, it should return ErrCommitNotFound.
	Checkout(dir string, timestamp string, repo *v1.Repository) (*v1.Commit, error)
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
		dir, err := c.gh.Clone(repo)
		if err != nil {
			logger.Error("failed to clone", "error", err)
			return fmt.Errorf("failed to clone %v: %w", repo, err)
		}
		slog.Info("cloned", "dir", dir)
		for _, ts := range c.config.GetTimestamps() {
			commit, err := c.gh.Checkout(dir, ts, repo)
			if errors.Is(err, ErrCommitNotFound) {
				logger.Warn("commit not found", "timestamp", ts)
				continue
			}
			if err != nil {
				logger.Error("failed to checkout", "error", err)
				return fmt.Errorf("failed to checkout %v: %w", ts, err)
			}
			slog.Info("checked out", "commit", commit)
		}
	}
	return nil
}

package gh

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nokamoto/covalyzer-go/internal/infra/log"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

type GitHub struct {
	wd string
}

func NewGitHub() (*GitHub, error) {
	wd, err := os.MkdirTemp("", "covalyzer")
	if err != nil {
		return nil, fmt.Errorf("failed to create a temporary directory: %w", err)
	}
	return &GitHub{
		wd: wd,
	}, nil
}

func (g *GitHub) Clone(repo *v1.Repository) (string, error) {
	github := repo.GetGh()
	if github != "" {
		github = fmt.Sprintf("%s/", github)
	}
	cmd := exec.Command("gh", "repo", "clone", fmt.Sprintf("%s%s/%s", github, repo.GetOwner(), repo.GetRepo()))
	cmd.Stdout = log.NewLogWriter()
	cmd.Stderr = log.NewErrorLogWriter()
	cmd.Dir = g.wd
	if err := cmd.Run(); err != nil {
		slog.Error("failed to clone a repository", "error", err)
		return "", fmt.Errorf("failed to clone a repository: %w", err)
	}
	return filepath.Join(g.wd, repo.GetRepo()), nil
}

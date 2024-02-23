package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nokamoto/covalyzer-go/internal/infra/command"
	"github.com/nokamoto/covalyzer-go/internal/usecase"
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
	arg := fmt.Sprintf("%s%s/%s", github, repo.GetOwner(), repo.GetRepo())
	if err := command.Run("gh", "repo", "clone", arg)(command.WithDir(g.wd)); err != nil {
		return "", err
	}
	return filepath.Join(g.wd, repo.GetRepo()), nil
}

func (g *GitHub) recentCommit(dir, timestamp string, repo *v1.Repository) (*v1.Commit, error) {
	var res bytes.Buffer
	api := fmt.Sprintf("/repos/%s/%s/commits?per_page=1&until=%s", repo.GetOwner(), repo.GetRepo(), timestamp)
	if err := command.Run("gh", "api", api)(command.WithDir(dir), command.WithStdout(&res)); err != nil {
		return nil, err
	}
	type commit struct {
		Sha string
	}
	var commits []commit
	if err := json.Unmarshal(res.Bytes(), &commits); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %v: %w", res.String(), err)
	}
	switch len(commits) {
	case 0:
		return nil, fmt.Errorf("%w: %v", usecase.ErrCommitNotFound, timestamp)
	case 1:
		return &v1.Commit{
			Sha: commits[0].Sha,
		}, nil
	}
	return nil, fmt.Errorf("unexpected commits: %v", commits)
}

func (g *GitHub) Checkout(dir, timestamp string, repo *v1.Repository) (*v1.Commit, error) {
	commit, err := g.recentCommit(dir, timestamp, repo)
	if err != nil {
		return nil, err
	}
	wd := fmt.Sprintf("%s/%s", g.wd, repo.GetRepo())
	if err := command.Run("git", "checkout", commit.GetSha())(command.WithDir(wd)); err != nil {
		return nil, err
	}
	return commit, nil
}

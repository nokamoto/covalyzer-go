package gh

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/nokamoto/covalyzer-go/internal/infra/command"
	"github.com/nokamoto/covalyzer-go/internal/usecase"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

type GitHub struct {
	wd command.WorkingDir
}

func NewGitHub(wd command.WorkingDir) (*GitHub, error) {
	return &GitHub{
		wd: wd,
	}, nil
}

func (g *GitHub) Clone(repo *v1.Repository) error {
	github := repo.GetGh()
	if github != "" {
		github = fmt.Sprintf("%s/", github)
	}
	arg := fmt.Sprintf("%s%s/%s", github, repo.GetOwner(), repo.GetRepo())
	if err := command.Run("gh", "repo", "clone", arg)(g.wd.WithDir()); err != nil {
		return err
	}
	return nil
}

func (g *GitHub) recentCommit(repo *v1.Repository, timestamp string) (*v1.Commit, error) {
	var res bytes.Buffer
	api := fmt.Sprintf("/repos/%s/%s/commits?per_page=1&until=%s", repo.GetOwner(), repo.GetRepo(), timestamp)
	if err := command.Run("gh", "api", api)(g.wd.WithRepoDir(repo), command.WithStdout(&res)); err != nil {
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

func (g *GitHub) Checkout(repo *v1.Repository, timestamp string) (*v1.Commit, error) {
	commit, err := g.recentCommit(repo, timestamp)
	if err != nil {
		return nil, err
	}
	if err := command.Run("git", "checkout", commit.GetSha())(g.wd.WithRepoDir(repo)); err != nil {
		return nil, err
	}
	return commit, nil
}

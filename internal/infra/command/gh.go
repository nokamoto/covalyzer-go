package command

import (
	"encoding/json"
	"fmt"

	"github.com/nokamoto/covalyzer-go/internal/usecase"
	"github.com/nokamoto/covalyzer-go/internal/util/xslices"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

type GitHub struct {
	wd     WorkingDir
	runner runner
}

func NewGitHub(wd WorkingDir) (*GitHub, error) {
	return &GitHub{
		wd:     wd,
		runner: &command{},
	}, nil
}

func (g *GitHub) Clone(repo *v1.Repository) error {
	github := repo.GetGh()
	if github != "" {
		github = fmt.Sprintf("%s/", github)
	}
	arg := fmt.Sprintf("%s%s/%s", github, repo.GetOwner(), repo.GetRepo())
	if err := g.runner.run("gh", xslices.Concat("repo", "clone", arg), g.wd.withDir()); err != nil {
		return err
	}
	return nil
}

func (g *GitHub) recentCommit(repo *v1.Repository, timestamp string) (*v1.Commit, error) {
	var opts []option
	opts = append(opts, g.wd.withRepoDir(repo))
	if repo.GetGh() != "" {
		env := map[string]string{}
		env["GH_HOST"] = repo.GetGh()
		opts = append(opts, withEnv(env))
	}
	api := fmt.Sprintf("/repos/%s/%s/commits?per_page=1&until=%s", repo.GetOwner(), repo.GetRepo(), timestamp)
	res, err := g.runner.runO("gh", xslices.Concat("api", api), opts...)
	if err != nil {
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
	if err := g.runner.run("git", xslices.Concat("checkout", commit.GetSha()), g.wd.withRepoDir(repo)); err != nil {
		return nil, err
	}
	return commit, nil
}

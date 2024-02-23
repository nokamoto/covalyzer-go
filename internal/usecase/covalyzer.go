//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock_test.go -package=$GOPACKAGE
package usecase

import (
	"errors"
	"fmt"

	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

var (
	ErrCommitNotFound = fmt.Errorf("commit not found")
)

type gh interface {
	// Clone clones a repository.
	Clone(*v1.Repository) error
	// Checkout checks out a repository at a specific timestamp and returns the commit.
	// The commit is the most recent commit before the timestamp.
	//
	// If the commit is not found for the timestamp, it should return ErrCommitNotFound.
	Checkout(repo *v1.Repository, timestamp string) (*v1.Commit, error)
}

type gotool interface {
	// Cover tests a repository and returns the coverage.
	Cover(repo *v1.Repository) (*v1.Cover, error)
}

type Covalyzer struct {
	config *v1.Config
	gh     gh
	gotool gotool
}

func NewCovalyzer(config *v1.Config, gh gh, gotool gotool) *Covalyzer {
	return &Covalyzer{
		config: config,
		gh:     gh,
		gotool: gotool,
	}
}

func (c *Covalyzer) Run() (*v1.Covalyzer, error) {
	var res v1.Covalyzer
	for _, repo := range c.config.GetRepositories() {
		rc := v1.RepositoryCoverages{
			Repository: repo,
		}
		err := c.gh.Clone(repo)
		if err != nil {
			return nil, fmt.Errorf("failed to clone %v: %w", repo, err)
		}

		for _, ts := range c.config.GetTimestamps() {
			commit, err := c.gh.Checkout(repo, ts)
			if errors.Is(err, ErrCommitNotFound) {
				rc.Coverages = append(rc.Coverages, &v1.Coverage{
					Commit: commit,
				})
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("failed to checkout %v: %w", ts, err)
			}

			cover, err := c.gotool.Cover(repo)
			if err != nil {
				return nil, fmt.Errorf("failed to test %v: %w", commit, err)
			}

			rc.Coverages = append(rc.Coverages, &v1.Coverage{
				Commit: commit,
				Cover:  cover,
			})
		}
		res.Repositories = append(res.Repositories, &rc)
	}
	return &res, nil
}

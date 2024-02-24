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
	CoverTotal(repo *v1.Repository) (float32, error)
	CoverGinkgoReport(repo *v1.Repository) ([]*v1.GinkgoReportCover, error)
	CoverGinkgoOutline(repo *v1.Repository) ([]*v1.GinkgoOutlineCover, error)
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

			total, err := c.gotool.CoverTotal(repo)
			if err != nil {
				return nil, fmt.Errorf("failed to analyze go tool cover %v: %w", commit, err)
			}

			outline, err := c.gotool.CoverGinkgoOutline(repo)
			if err != nil {
				return nil, fmt.Errorf("failed to analyze ginkgo outline %v: %w", commit, err)
			}

			report, err := c.gotool.CoverGinkgoReport(repo)
			if err != nil {
				return nil, fmt.Errorf("failed to analyze ginkgo run %v: %w", commit, err)
			}

			rc.Coverages = append(rc.Coverages, &v1.Coverage{
				Commit: commit,
				Cover: &v1.Cover{
					Total:          total,
					GinkgoOutlines: outline,
					GinkgoReports:  report,
				},
			})
		}
		res.Repositories = append(res.Repositories, &rc)
	}
	return &res, nil
}

package gh

import v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"

type GitHub struct{}

func NewGitHub() *GitHub {
	return &GitHub{}
}

func (g *GitHub) Clone(repo *v1.Repository) error {
	return nil
}

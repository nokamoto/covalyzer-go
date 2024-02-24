package command

import (
	"fmt"
	"os"
	"path/filepath"

	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

type WorkingDir string

// NewWorkingDir creates a temporary directory.
func NewWorkingDir() (WorkingDir, error) {
	wd, err := os.MkdirTemp("", "covalyzer")
	if err != nil {
		return "", fmt.Errorf("failed to create a temporary directory: %w", err)
	}
	return WorkingDir(wd), nil
}

func (w WorkingDir) withDir() option {
	return withDir(string(w))
}

func (w WorkingDir) withRepoDir(repo *v1.Repository) option {
	return withDir(filepath.Join(string(w), repo.GetRepo()))
}

// Clean removes the working directory.
func (w WorkingDir) Clean() {
	os.RemoveAll(string(w))
}

package command

import (
	"fmt"
	"log/slog"
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
	slog.Debug("created a working directory", "dir", wd)
	return WorkingDir(wd), nil
}

// WithDir sets the working directory for the command.
func (w WorkingDir) WithDir() Option {
	return WithDir(string(w))
}

// WithRepoDir sets the working directory for the command to the repository directory.
func (w WorkingDir) WithRepoDir(repo *v1.Repository) Option {
	return WithDir(filepath.Join(string(w), repo.GetRepo()))
}

// Clean removes the working directory.
func (w WorkingDir) Clean() {
	os.RemoveAll(string(w))
}

package command

import (
	"fmt"
	"os/exec"

	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	"go.uber.org/mock/gomock"
)

type withDirMatcher struct {
	f option
}

func newWithDirMatcher(wd WorkingDir) gomock.Matcher {
	return &withDirMatcher{f: wd.withDir()}
}

func newWithRepoDirMatcher(wd WorkingDir, repo *v1.Repository) gomock.Matcher {
	return &withDirMatcher{f: wd.withRepoDir(repo)}
}

func (m *withDirMatcher) Matches(x any) bool {
	f, ok := x.(option)
	if !ok {
		return false
	}
	var c1, c2 exec.Cmd
	m.f(&c1)
	f(&c2)
	return c1.Dir == c2.Dir
}

func (m *withDirMatcher) String() string {
	var c1 exec.Cmd
	m.f(&c1)
	return fmt.Sprintf("option(dir=%s)", c1.Dir)
}

package gotool

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/nokamoto/covalyzer-go/internal/infra/command"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

type GoTool struct {
	wd command.WorkingDir
}

func NewGoTool(wd command.WorkingDir) *GoTool {
	return &GoTool{
		wd: wd,
	}
}

func (g *GoTool) Cover(repo *v1.Repository) (*v1.Cover, error) {
	const out = "c.out"
	if err := command.Run("go", "test", "-coverprofile", out, "./...")(g.wd.WithRepoDir(repo)); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := command.Run("go", "tool", "cover", "-func", out)(g.wd.WithRepoDir(repo), command.WithStdout(&buf)); err != nil {
		return nil, err
	}
	// parse last line of the output
	// total: (statements) 82.2%
	var line string
	scan := bufio.NewScanner(&buf)
	for scan.Scan() {
		line = scan.Text()
	}
	if err := scan.Err(); err != nil {
		slog.Debug("failed to scan", "error", err)
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	ss := strings.Fields(line)
	if len(ss) == 0 {
		return nil, fmt.Errorf("unexpected output: %v", line)
	}
	last := ss[len(ss)-1]
	last = strings.TrimSuffix(last, "%")
	total, err := strconv.ParseFloat(last, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}
	return &v1.Cover{
		Total: float32(total),
	}, nil
}

package gotool

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/nokamoto/covalyzer-go/internal/infra/command"
	"github.com/nokamoto/covalyzer-go/internal/util/xslices"
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

func (g *GoTool) testPackages(repo *v1.Repository) ([]string, error) {
	var buf bytes.Buffer
	if err := command.Run("go", "list", "-f", "{{.Dir}}", "./...")(g.wd.WithRepoDir(repo), command.WithStdout(&buf)); err != nil {
		return nil, err
	}
	var pkgs []string
	scan := bufio.NewScanner(&buf)
	for scan.Scan() {
		pkgs = append(pkgs, scan.Text())
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	pkgs = slices.DeleteFunc(pkgs, func(s string) bool {
		for _, ginkgo := range repo.GetGinkgoPackages() {
			if strings.Contains(s, ginkgo) {
				return true
			}
		}
		return false
	})
	return pkgs, nil
}

func (g *GoTool) parseTotal(buf bytes.Buffer) (float32, error) {
	var line string
	scan := bufio.NewScanner(&buf)
	for scan.Scan() {
		line = scan.Text()
	}
	if err := scan.Err(); err != nil {
		slog.Debug("failed to scan", "error", err)
		return 0, fmt.Errorf("failed to scan: %w", err)
	}

	ss := strings.Fields(line)
	if len(ss) == 0 {
		return 0, fmt.Errorf("unexpected output: %v", line)
	}
	last := ss[len(ss)-1]
	last = strings.TrimSuffix(last, "%")
	total, err := strconv.ParseFloat(last, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse: %w", err)
	}
	return float32(total), nil
}

func (g *GoTool) ginkgoCover(repo *v1.Repository) ([]*v1.GinkgoCover, error) {
	var res []*v1.GinkgoCover
	dir := filepath.Join(string(g.wd), repo.GetRepo())
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") && !info.IsDir() {
			var buf bytes.Buffer
			if err := command.Run("ginkgo", "outline", "--format", "json", path)(command.WithStdout(&buf)); err != nil {
				return nil
			}

			cover := v1.GinkgoCover{
				File: strings.TrimPrefix(path, dir+"/"),
			}

			// https://github.com/onsi/ginkgo/blob/cd418b74c1e8502b305aab84e882516a6335a0e7/ginkgo/outline/outline.go#L70-L72
			type metadata struct {
				Name string
				Spec bool
			}
			type node struct {
				metadata
				Nodes []*node
			}
			var o []*node
			if err := json.Unmarshal(buf.Bytes(), &o); err != nil {
				return fmt.Errorf("failed to unmarshal: %s: %w", path, err)
			}

			var rec func(*node)
			rec = func(n *node) {
				if n.Spec {
					cover.OutlineNodes++
				}
				for _, c := range n.Nodes {
					rec(c)
				}
			}
			for _, n := range o {
				rec(n)
			}

			res = append(res, &cover)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, err
}

func (g *GoTool) Cover(repo *v1.Repository) (*v1.Cover, error) {
	list, err := g.testPackages(repo)
	if err != nil {
		return nil, err
	}

	const out = "c.out"
	if err := command.Run("go", xslices.Concat("test", "-coverprofile", out, list)...)(g.wd.WithRepoDir(repo)); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := command.Run("go", "tool", "cover", "-func", out)(g.wd.WithRepoDir(repo), command.WithStdout(&buf)); err != nil {
		return nil, err
	}

	total, err := g.parseTotal(buf)
	if err != nil {
		return nil, err
	}

	ginkgo, err := g.ginkgoCover(repo)
	if err != nil {
		return nil, err
	}

	return &v1.Cover{
		Total:  float32(total),
		Ginkgo: ginkgo,
	}, nil
}

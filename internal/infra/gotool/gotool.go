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
	"github.com/onsi/ginkgo/v2/types"
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

func (g *GoTool) ginkgoOutlineCover(repo *v1.Repository) ([]*v1.GinkgoOutlineCover, error) {
	var res []*v1.GinkgoOutlineCover
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

			cover := v1.GinkgoOutlineCover{
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

func (g *GoTool) ginkgoCover(repo *v1.Repository) ([]*v1.GinkgoReportCover, error) {
	const out = "report.json"
	var res []*v1.GinkgoReportCover
	var buf bytes.Buffer
	for _, pkg := range repo.GetGinkgoPackages() {
		if err := command.Run("ginkgo", "run", "--dry-run", fmt.Sprintf("--json-report=%s", out), pkg)(g.wd.WithRepoDir(repo), command.WithStdout(&buf)); err != nil {
			return nil, err
		}

		file := filepath.Join(string(g.wd), repo.GetRepo(), out)
		bs, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read: %s: %w", file, err)
		}

		var report []types.Report
		if err := json.Unmarshal(bs, &report); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s: %w", file, err)
		}

		var suites []*v1.GinkgoSuiteCover
		for _, r := range report {
			suites = append(suites, &v1.GinkgoSuiteCover{
				Description:      r.SuiteDescription,
				TotalSpecs:       int32(r.PreRunStats.TotalSpecs),
				SpecsThatWillRun: int32(r.PreRunStats.SpecsThatWillRun),
			})
		}

		res = append(res, &v1.GinkgoReportCover{
			Package: pkg,
			Suites:  suites,
		})
	}
	return res, nil
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

	outline, err := g.ginkgoOutlineCover(repo)
	if err != nil {
		return nil, err
	}

	report, err := g.ginkgoCover(repo)
	if err != nil {
		return nil, err
	}

	return &v1.Cover{
		Total:          float32(total),
		GinkgoOutlines: outline,
		GinkgoReports:  report,
	}, nil
}

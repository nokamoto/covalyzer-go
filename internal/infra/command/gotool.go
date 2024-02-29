package command

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

	"github.com/google/uuid"
	"github.com/nokamoto/covalyzer-go/internal/util/xslices"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	"github.com/onsi/ginkgo/v2/types"
)

type GoTool struct {
	wd     WorkingDir
	runner runner
	random func() string
}

func NewGoTool(wd WorkingDir) *GoTool {
	return &GoTool{
		wd:     wd,
		runner: &command{},
		random: uuid.NewString,
	}
}

func (g *GoTool) testPackages(repo *v1.Repository) ([]string, error) {
	buf, err := g.runner.runO("go", xslices.Concat("list", "-f", "{{.Dir}}", "./..."), g.wd.withRepoDir(repo))
	if err != nil {
		// ignore error because go list may fail if the repository is not a go project
		return nil, nil
	}
	var pkgs []string
	scan := bufio.NewScanner(buf)
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

func (g *GoTool) parseTotal(buf *bytes.Buffer) (float32, error) {
	var line string
	scan := bufio.NewScanner(buf)
	for scan.Scan() {
		line = scan.Text()
	}
	if err := scan.Err(); err != nil {
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

var ginkgo = []string{"run", "github.com/onsi/ginkgo/v2/ginkgo"}

func (g *GoTool) CoverGinkgoOutline(repo *v1.Repository) ([]*v1.GinkgoOutlineCover, error) {
	var res []*v1.GinkgoOutlineCover
	dir := filepath.Join(string(g.wd), repo.GetRepo())
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") && !info.IsDir() {
			buf, err := g.runner.runO("go", xslices.Concat(ginkgo, "outline", "--format", "json", path))
			if err != nil {
				// ignore error because ginkgo outline may fail if there is no test
				return nil
			}

			cover := v1.GinkgoOutlineCover{
				File: strings.TrimPrefix(path, dir+"/"),
			}

			// https://github.com/onsi/ginkgo/blob/cd418b74c1e8502b305aab84e882516a6335a0e7/ginkgo/outline/outline.go#L70-L72
			type metadata struct {
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

func (g *GoTool) CoverGinkgoReport(repo *v1.Repository) ([]*v1.GinkgoReportCover, error) {
	const out = "report.json"
	var res []*v1.GinkgoReportCover
	for _, pkg := range repo.GetGinkgoPackages() {
		err := g.runner.run("go", xslices.Concat(ginkgo, "run", "--dry-run", fmt.Sprintf("--json-report=%s", out), pkg), g.wd.withRepoDir(repo))
		if err != nil {
			// ignore error because ginkgo run may fail if there is no test
			res = append(res, &v1.GinkgoReportCover{
				Package: pkg,
			})
			continue
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

var (
	errCoverProfileNotFound = fmt.Errorf("failed to write cover profile")
)

func (g *GoTool) CoverTotal(repo *v1.Repository) (float32, error) {
	list, err := g.testPackages(repo)
	if err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, nil
	}

	out := "c.out"
	if g.random != nil {
		out = fmt.Sprintf("%s.%s", out, g.random())
	}
	if err := g.runner.run("go", xslices.Concat("test", "-coverprofile", out, list), g.wd.withRepoDir(repo)); err != nil {
		slog.Debug("[ignored] failed to run go test -coverprofile", "profile", out, "err", err)
	}

	if _, err := os.Stat(filepath.Join(string(g.wd), repo.GetRepo(), out)); err != nil {
		return 0, fmt.Errorf("%w: %w", errCoverProfileNotFound, err)
	}

	buf, err := g.runner.runO("go", xslices.Concat("tool", "cover", "-func", out), g.wd.withRepoDir(repo))
	if err != nil {
		return 0, err
	}

	return g.parseTotal(buf)
}

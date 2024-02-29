package command

import (
	"bytes"
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"
)

//go:embed testdata/empty-report.json
var emptyReportJSON []byte

//go:embed testdata/report.json
var reportJSON []byte

//go:embed testdata/outline.json
var outlineJSON string

func TestGoTool_CoverTotal(t *testing.T) {
	internalErr := errors.New("internal")

	tests := []struct {
		name       string
		repo       *v1.Repository
		filesystem func(dir string)
		mock       func(*Mockrunner, WorkingDir)
		want       float32
		wantErr    error
	}{
		{
			name: "ok",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			filesystem: func(dir string) {
				_ = os.WriteFile(filepath.Join(dir, "c.out"), []byte{}, 0644)
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				opt := newWithRepoDirMatcher(wd, &v1.Repository{
					Owner: "foo",
					Repo:  "bar",
				})
				m.EXPECT().runO(
					"go",
					[]string{"list", "-f", "{{.Dir}}", "./..."},
					opt,
				).Return(
					bytes.NewBufferString("package1\npackage2\n"),
					nil,
				)
				m.EXPECT().run(
					"go",
					[]string{"test", "-coverprofile", "c.out", "package1", "package2"},
					opt,
				).Return(nil)
				m.EXPECT().runO(
					"go",
					[]string{"tool", "cover", "-func", "c.out"},
					opt,
				).Return(
					bytes.NewBufferString("...\nxxx xxx 12.3%\n"),
					nil,
				)
			},
			want: 12.3,
		},
		{
			name: "ok with ginkgo packages",
			repo: &v1.Repository{
				Owner:          "foo",
				Repo:           "bar",
				GinkgoPackages: []string{"package1"},
			},
			filesystem: func(dir string) {
				_ = os.WriteFile(filepath.Join(dir, "c.out"), []byte{}, 0644)
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				m.EXPECT().runO(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					bytes.NewBufferString("package1\npackage2\n"),
					nil,
				)
				m.EXPECT().run(
					gomock.Any(),
					[]string{"test", "-coverprofile", "c.out", "package2"},
					gomock.Any(),
				).Return(nil)
				m.EXPECT().runO(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					bytes.NewBufferString("...\nxxx xxx 12.3%\n"),
					nil,
				)
			},
			want: 12.3,
		},
		{
			name: "return 0 if go list failed",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				m.EXPECT().runO(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, internalErr)
			},
		},
		{
			name: "continue if go test failed",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			filesystem: func(dir string) {
				_ = os.WriteFile(filepath.Join(dir, "c.out"), []byte{}, 0644)
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				m.EXPECT().runO(gomock.Any(), gomock.Any(), gomock.Any()).Return(bytes.NewBufferString("a"), nil)
				m.EXPECT().run(gomock.Any(), gomock.Any(), gomock.Any()).Return(internalErr)
				m.EXPECT().runO(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					bytes.NewBufferString("...\nxxx xxx 12.3%\n"),
					nil,
				)
			},
			want: 12.3,
		},
		{
			name: "error if cover profile not found",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				m.EXPECT().runO(gomock.Any(), gomock.Any(), gomock.Any()).Return(bytes.NewBufferString("a"), nil)
				m.EXPECT().run(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: errCoverProfileNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd := WorkingDir(t.TempDir())
			ctrl := gomock.NewController(t)
			m := NewMockrunner(ctrl)
			g := &GoTool{
				wd:     wd,
				runner: m,
			}

			repoDir := filepath.Join(string(wd), tt.repo.GetRepo())
			_ = os.MkdirAll(repoDir, 0755)
			if tt.filesystem != nil {
				tt.filesystem(repoDir)
			}

			if tt.mock != nil {
				tt.mock(m, wd)
			}
			got, err := g.CoverTotal(tt.repo)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GoTool.CoverTotal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GoTool.CoverTotal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGoTool_CoverGinkgoOutline(t *testing.T) {
	tests := []struct {
		name       string
		repo       *v1.Repository
		filesystem func(WorkingDir)
		mock       func(*Mockrunner, WorkingDir)
		want       []*v1.GinkgoOutlineCover
		wantErr    error
	}{
		{
			name: "ok with empty directory",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
		},
		{
			name: "ok with go files",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			filesystem: func(wd WorkingDir) {
				var empty []byte
				dir := filepath.Join(string(wd), "bar")
				_ = os.MkdirAll(filepath.Join(dir, "dir1"), 0755)
				_ = os.WriteFile(filepath.Join(dir, "main.go"), empty, 0644)
				_ = os.WriteFile(filepath.Join(dir, "dir1", "main.go"), empty, 0644)
				_ = os.WriteFile(filepath.Join(dir, "ignore.json"), empty, 0644)
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				m.EXPECT().runO(
					"ginkgo",
					[]string{"outline", "--format", "json", filepath.Join(string(wd), "bar/main.go")},
				).Return(
					bytes.NewBufferString("[]"),
					nil,
				)
				m.EXPECT().runO(
					"ginkgo",
					[]string{"outline", "--format", "json", filepath.Join(string(wd), "bar/dir1/main.go")},
				).Return(
					bytes.NewBufferString(outlineJSON),
					nil,
				)
			},
			want: []*v1.GinkgoOutlineCover{
				{
					File:         "dir1/main.go",
					OutlineNodes: 2,
				},
				{
					File: "main.go",
				},
			},
		},
		{
			name: "continue on error",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			filesystem: func(wd WorkingDir) {
				var empty []byte
				dir := filepath.Join(string(wd), "bar")
				_ = os.WriteFile(filepath.Join(dir, "main.go"), empty, 0644)
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				m.EXPECT().runO(
					"ginkgo",
					[]string{"outline", "--format", "json", filepath.Join(string(wd), "bar/main.go")},
				).Return(
					nil,
					errors.New("internal"),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd := WorkingDir(t.TempDir())
			ctrl := gomock.NewController(t)
			m := NewMockrunner(ctrl)
			g := &GoTool{
				wd:     wd,
				runner: m,
			}
			repoDir := filepath.Join(string(wd), tt.repo.GetRepo())
			_ = os.MkdirAll(repoDir, 0755)
			if tt.filesystem != nil {
				tt.filesystem(wd)
			}
			if tt.mock != nil {
				tt.mock(m, wd)
			}
			got, err := g.CoverGinkgoOutline(tt.repo)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GoTool.CoverGinkgoOutline() error = %v, wantErr %v", err, tt.wantErr)
			}
			opts := cmp.Options{
				// sort for walk order
				protocmp.SortRepeated(func(a, b *v1.GinkgoOutlineCover) bool {
					return a.File < b.File
				}),
				protocmp.Transform(),
			}
			if diff := cmp.Diff(tt.want, got, opts); diff != "" {
				t.Errorf("GoTool.CoverGinkgoOutline() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGoTool_CoverGinkgoReport(t *testing.T) {
	internalErr := errors.New("internal")

	tests := []struct {
		name       string
		repo       *v1.Repository
		filesystem func(WorkingDir)
		mock       func(*Mockrunner, WorkingDir)
		want       []*v1.GinkgoReportCover
		wantErr    error
	}{
		{
			name: "ok with empty report",
			repo: &v1.Repository{
				Owner:          "foo",
				Repo:           "bar",
				GinkgoPackages: []string{"package1"},
			},
			filesystem: func(wd WorkingDir) {
				_ = os.WriteFile(filepath.Join(string(wd), "bar", "report.json"), emptyReportJSON, 0644)
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				opt := newWithRepoDirMatcher(wd, &v1.Repository{
					Owner: "foo",
					Repo:  "bar",
				})
				m.EXPECT().run(
					"ginkgo",
					[]string{"run", "--dry-run", "--json-report=report.json", "package1"},
					opt,
				).Return(
					nil,
				)
			},
			want: []*v1.GinkgoReportCover{
				{
					Package: "package1",
				},
			},
		},
		{
			name: "ok with report",
			repo: &v1.Repository{
				Owner:          "foo",
				Repo:           "bar",
				GinkgoPackages: []string{"package1"},
			},
			filesystem: func(wd WorkingDir) {
				_ = os.WriteFile(filepath.Join(string(wd), "bar", "report.json"), reportJSON, 0644)
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				m.EXPECT().run(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: []*v1.GinkgoReportCover{
				{
					Package: "package1",
					Suites: []*v1.GinkgoSuiteCover{
						{
							Description:      "suite1",
							TotalSpecs:       1,
							SpecsThatWillRun: 2,
						},
						{
							Description:      "suite2",
							TotalSpecs:       3,
							SpecsThatWillRun: 4,
						},
					},
				},
			},
		},
		{
			name: "continue even if ginkgo run failed",
			repo: &v1.Repository{
				Owner:          "foo",
				Repo:           "bar",
				GinkgoPackages: []string{"package1", "package2"},
			},
			filesystem: func(wd WorkingDir) {
				_ = os.WriteFile(filepath.Join(string(wd), "bar", "report.json"), reportJSON, 0644)
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				gomock.InOrder(
					m.EXPECT().run(gomock.Any(), gomock.Any(), gomock.Any()).Return(internalErr),
					m.EXPECT().run(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			want: []*v1.GinkgoReportCover{
				{
					Package: "package1",
				},
				{
					Package: "package2",
					Suites: []*v1.GinkgoSuiteCover{
						{
							Description:      "suite1",
							TotalSpecs:       1,
							SpecsThatWillRun: 2,
						},
						{
							Description:      "suite2",
							TotalSpecs:       3,
							SpecsThatWillRun: 4,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd := WorkingDir(t.TempDir())
			ctrl := gomock.NewController(t)
			m := NewMockrunner(ctrl)
			g := &GoTool{
				wd:     wd,
				runner: m,
			}
			repoDir := filepath.Join(string(wd), tt.repo.GetRepo())
			_ = os.MkdirAll(repoDir, 0755)
			if tt.filesystem != nil {
				tt.filesystem(wd)
			}
			if tt.mock != nil {
				tt.mock(m, wd)
			}
			got, err := g.CoverGinkgoReport(tt.repo)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GoTool.CoverGinkgoReport() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GoTool.CoverGinkgoReport() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

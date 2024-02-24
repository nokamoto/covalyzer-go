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

func TestGoTool_Cover(t *testing.T) {
	internalErr := errors.New("internal")

	tests := []struct {
		name       string
		repo       *v1.Repository
		filesystem func(WorkingDir)
		mock       func(*Mockrunner, WorkingDir)
		want       *v1.Cover
		wantErr    error
	}{
		{
			name: "ok",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			filesystem: func(wd WorkingDir) {
				repoDir := filepath.Join(string(wd), "bar")
				_ = os.MkdirAll(repoDir, 0755)
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
			want: &v1.Cover{
				Total: 12.3,
			},
		},
		{
			name: "ok with ginkgo packages",
			repo: &v1.Repository{
				Owner:          "foo",
				Repo:           "bar",
				GinkgoPackages: []string{"package1"},
			},
			filesystem: func(wd WorkingDir) {
				repoDir := filepath.Join(string(wd), "bar")
				_ = os.MkdirAll(repoDir, 0755)
				_ = os.WriteFile(filepath.Join(repoDir, "report.json"), emptyReportJSON, 0644)
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				opt := newWithRepoDirMatcher(wd, &v1.Repository{
					Owner: "foo",
					Repo:  "bar",
				})
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
				m.EXPECT().run(
					"ginkgo",
					[]string{"run", "--dry-run", "--json-report=report.json", "package1"},
					opt,
				).Return(nil)
			},
			want: &v1.Cover{
				Total: 12.3,
				GinkgoReports: []*v1.GinkgoReportCover{
					{
						Package: "package1",
					},
				},
			},
		},
		{
			name: "failed to go list",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			mock: func(m *Mockrunner, wd WorkingDir) {
				m.EXPECT().runO(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, internalErr)
			},
			wantErr: internalErr,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			wd := WorkingDir(t.TempDir())
			ctrl := gomock.NewController(t)
			m := NewMockrunner(ctrl)
			g := &GoTool{
				wd:     wd,
				runner: m,
			}
			if tt.filesystem != nil {
				tt.filesystem(wd)
			}
			if tt.mock != nil {
				tt.mock(m, wd)
			}
			got, err := g.Cover(tt.repo)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GoTool.Cover() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GoTool.Cover() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

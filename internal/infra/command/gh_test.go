package command

import (
	"bytes"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/covalyzer-go/internal/usecase"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGitHub_Clone(t *testing.T) {
	wd := WorkingDir(t.TempDir())

	tests := []struct {
		name    string
		repo    *v1.Repository
		mock    func(*Mockrunner)
		wantErr bool
	}{
		{
			name: "ok",
			repo: &v1.Repository{
				Owner: "bar",
				Repo:  "foo",
			},
			mock: func(m *Mockrunner) {
				m.EXPECT().run(
					"gh",
					[]string{"repo", "clone", "bar/foo"},
					newWithDirMatcher(wd),
				).Return(nil)
			},
		},
		{
			name: "ok with ghe",
			repo: &v1.Repository{
				Gh:    "example.com",
				Owner: "bar",
				Repo:  "foo",
			},
			mock: func(m *Mockrunner) {
				m.EXPECT().run(
					"gh",
					[]string{"repo", "clone", "example.com/bar/foo"},
					newWithDirMatcher(wd),
				).Return(nil)
			},
		},
		{
			name: "error",
			mock: func(m *Mockrunner) {
				m.EXPECT().run(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := NewMockrunner(ctrl)
			if tt.mock != nil {
				tt.mock(m)
			}
			g := &GitHub{
				wd:     wd,
				runner: m,
			}
			if err := g.Clone(tt.repo); (err != nil) != tt.wantErr {
				t.Errorf("GitHub.Clone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitHub_Checkout(t *testing.T) {
	wd := WorkingDir(t.TempDir())
	ts := "2024-02-01T00:00:00Z"
	internalErr := errors.New("internal error")

	tests := []struct {
		name    string
		repo    *v1.Repository
		mock    func(*Mockrunner)
		want    *v1.Commit
		wantErr error
	}{
		{
			name: "ok",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			mock: func(m *Mockrunner) {
				m.EXPECT().runO(
					"gh",
					[]string{"api", "/repos/foo/bar/commits?per_page=1&until=2024-02-01T00:00:00Z"},
					newWithRepoDirMatcher(wd, &v1.Repository{
						Owner: "foo",
						Repo:  "bar",
					}),
				).Return(bytes.NewBufferString(`[{"sha":"0"}]`), nil)
				m.EXPECT().run(
					"git",
					[]string{"checkout", "0"},
					newWithRepoDirMatcher(wd, &v1.Repository{
						Owner: "foo",
						Repo:  "bar",
					}),
				).Return(nil)
			},
			want: &v1.Commit{
				Sha: "0",
			},
		},
		{
			name: "ok with ghe",
			repo: &v1.Repository{
				Gh:    "example.com",
				Owner: "foo",
				Repo:  "bar",
			},
			mock: func(m *Mockrunner) {
				m.EXPECT().runO(
					"gh",
					[]string{"api", "/repos/foo/bar/commits?per_page=1&until=2024-02-01T00:00:00Z"},
					newWithRepoDirMatcher(wd, &v1.Repository{
						Owner: "foo",
						Repo:  "bar",
					}),
					newWithEnvMatcher(map[string]string{
						"GH_HOST": "example.com",
					}),
				).Return(bytes.NewBufferString(`[{"sha":"0"}]`), nil)
				m.EXPECT().run(
					"git",
					[]string{"checkout", "0"},
					newWithRepoDirMatcher(wd, &v1.Repository{
						Owner: "foo",
						Repo:  "bar",
					}),
				).Return(nil)
			},
			want: &v1.Commit{
				Sha: "0",
			},
		},
		{
			name: "failed to gh api",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			mock: func(m *Mockrunner) {
				m.EXPECT().runO(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, internalErr)
			},
			wantErr: internalErr,
		},
		{
			name: "ErrCommitNotFound if gh api returns empty",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			mock: func(m *Mockrunner) {
				m.EXPECT().runO(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(bytes.NewBufferString(`[]`), nil)
			},
			wantErr: usecase.ErrCommitNotFound,
		},
		{
			name: "failed to git checkout",
			repo: &v1.Repository{
				Owner: "foo",
				Repo:  "bar",
			},
			mock: func(m *Mockrunner) {
				m.EXPECT().runO(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(bytes.NewBufferString(`[{"sha":"0"}]`), nil)
				m.EXPECT().run(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(internalErr)
			},
			wantErr: internalErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := NewMockrunner(ctrl)
			if tt.mock != nil {
				tt.mock(m)
			}
			g := &GitHub{
				wd:     wd,
				runner: m,
			}
			got, err := g.Checkout(tt.repo, ts)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GitHub.Checkout() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GitHub.Checkout() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

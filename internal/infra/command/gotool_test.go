package command

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGoTool_Cover(t *testing.T) {
	repo := &v1.Repository{
		Owner: "foo",
		Repo:  "bar",
	}
	wd := WorkingDir(t.TempDir())
	internalErr := errors.New("internal")

	tests := []struct {
		name    string
		mock    func(*Mockrunner)
		want    *v1.Cover
		wantErr error
	}{
		{
			name: "ok",
			mock: func(m *Mockrunner) {},
			want: &v1.Cover{},
		},
		{
			name:    "failed to go list",
			mock:    func(m *Mockrunner) {},
			wantErr: internalErr,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := NewMockrunner(ctrl)
			g := &GoTool{
				wd:     wd,
				runner: m,
			}
			if tt.mock != nil {
				tt.mock(m)
			}
			got, err := g.Cover(repo)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GoTool.Cover() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GoTool.Cover() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

package usecase

import (
	"errors"
	"testing"

	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	gomock "go.uber.org/mock/gomock"
)

func TestCovalyzer_Run(t *testing.T) {
	internalErr := errors.New("internal error")
	config := &v1.Config{
		Repositories: []*v1.Repository{
			{
				Owner: "foo",
				Repo:  "bar",
			},
		},
	}
	tests := []struct {
		name    string
		mock    func(*Mockgh)
		wantErr error
	}{
		{
			name: "ok",
			mock: func(m *Mockgh) {
				m.EXPECT().Clone(config.GetRepositories()[0]).Return("", nil)
			},
		},
		{
			name: "failed to clone",
			mock: func(m *Mockgh) {
				m.EXPECT().Clone(gomock.Any()).Return("", internalErr)
			},
			wantErr: internalErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			gh := NewMockgh(ctrl)
			if tt.mock != nil {
				tt.mock(gh)
			}
			c := NewCovalyzer(config, gh)
			err := c.Run()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

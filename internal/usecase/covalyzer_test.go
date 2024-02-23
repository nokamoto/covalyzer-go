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
			{
				Owner: "baz",
				Repo:  "qux",
			},
		},
		Timestamps: []string{"0", "1"},
	}
	r0 := config.GetRepositories()[0]
	r1 := config.GetRepositories()[1]
	tests := []struct {
		name    string
		mock    func(*Mockgh)
		wantErr error
	}{
		{
			name: "ok",
			mock: func(m *Mockgh) {
				m.EXPECT().Clone(r0).Return("dir1", nil)
				m.EXPECT().Clone(r1).Return("dir2", nil)
				m.EXPECT().Checkout("dir1", "0", r0).Return(&v1.Commit{}, nil)
				m.EXPECT().Checkout("dir1", "1", r0).Return(&v1.Commit{}, nil)
				m.EXPECT().Checkout("dir2", "0", r1).Return(&v1.Commit{}, nil)
				m.EXPECT().Checkout("dir2", "1", r1).Return(&v1.Commit{}, nil)
			},
		},
		{
			name: "failed to clone",
			mock: func(m *Mockgh) {
				m.EXPECT().Clone(gomock.Any()).Return("", internalErr)
			},
			wantErr: internalErr,
		},
		{
			name: "failed to checkout",
			mock: func(m *Mockgh) {
				m.EXPECT().Clone(gomock.Any()).Return("", nil)
				m.EXPECT().Checkout(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, internalErr)
			},
			wantErr: internalErr,
		},
		{
			name: "continue on commit not found",
			mock: func(m *Mockgh) {
				m.EXPECT().Clone(r0).Return("dir1", nil)
				m.EXPECT().Clone(r1).Return("dir2", nil)
				m.EXPECT().Checkout("dir1", "0", r0).Return(&v1.Commit{}, nil)
				m.EXPECT().Checkout("dir1", "1", r0).Return(nil, ErrCommitNotFound)
				m.EXPECT().Checkout("dir2", "0", r1).Return(&v1.Commit{}, nil)
				m.EXPECT().Checkout("dir2", "1", r1).Return(&v1.Commit{}, nil)
			},
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

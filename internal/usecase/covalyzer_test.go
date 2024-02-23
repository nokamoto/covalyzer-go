package usecase

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/protobuf/testing/protocmp"
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
		mock    func(*Mockgh, *Mockgotool)
		want    *v1.Covalyzer
		wantErr error
	}{
		{
			name: "ok",
			mock: func(gh *Mockgh, gt *Mockgotool) {
				gh.EXPECT().Clone(r0).Return(nil)
				gh.EXPECT().Clone(r1).Return(nil)

				gh.EXPECT().Checkout(r0, "0").Return(&v1.Commit{
					Sha: "sha0",
				}, nil)
				gt.EXPECT().Cover(r0).Return(&v1.Cover{
					Total: 0.1,
				}, nil)

				gh.EXPECT().Checkout(r0, "1").Return(&v1.Commit{
					Sha: "sha1",
				}, nil)
				gt.EXPECT().Cover(r0).Return(&v1.Cover{
					Total: 0.2,
				}, nil)

				gh.EXPECT().Checkout(r1, "0").Return(&v1.Commit{
					Sha: "sha2",
				}, nil)
				gt.EXPECT().Cover(r1).Return(&v1.Cover{
					Total: 0.3,
				}, nil)

				gh.EXPECT().Checkout(r1, "1").Return(&v1.Commit{
					Sha: "sha3",
				}, nil)
				gt.EXPECT().Cover(r1).Return(&v1.Cover{
					Total: 0.4,
				}, nil)
			},
			want: &v1.Covalyzer{
				Repositories: []*v1.RepositoryCoverages{
					{
						Repository: r0,
						Coverages: []*v1.Coverage{
							{
								Commit: &v1.Commit{
									Sha: "sha0",
								},
								Cover: &v1.Cover{
									Total: 0.1,
								},
							},
							{
								Commit: &v1.Commit{
									Sha: "sha1",
								},
								Cover: &v1.Cover{
									Total: 0.2,
								},
							},
						},
					},
					{
						Repository: r1,
						Coverages: []*v1.Coverage{
							{
								Commit: &v1.Commit{
									Sha: "sha2",
								},
								Cover: &v1.Cover{
									Total: 0.3,
								},
							},
							{
								Commit: &v1.Commit{
									Sha: "sha3",
								},
								Cover: &v1.Cover{
									Total: 0.4,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "failed to clone",
			mock: func(gh *Mockgh, gt *Mockgotool) {
				gh.EXPECT().Clone(gomock.Any()).Return(internalErr)
			},
			wantErr: internalErr,
		},
		{
			name: "failed to checkout",
			mock: func(gh *Mockgh, gt *Mockgotool) {
				gh.EXPECT().Clone(gomock.Any()).Return(nil)
				gh.EXPECT().Checkout(gomock.Any(), gomock.Any()).Return(nil, internalErr)
			},
			wantErr: internalErr,
		},
		{
			name: "failed to test",
			mock: func(gh *Mockgh, gt *Mockgotool) {
				gh.EXPECT().Clone(gomock.Any()).Return(nil)
				gh.EXPECT().Checkout(gomock.Any(), gomock.Any()).Return(&v1.Commit{}, nil)
				gt.EXPECT().Cover(gomock.Any()).Return(nil, internalErr)
			},
			wantErr: internalErr,
		},
		{
			name: "continue on commit not found",
			mock: func(gh *Mockgh, gt *Mockgotool) {
				gh.EXPECT().Clone(r0).Return(nil)
				gh.EXPECT().Clone(r1).Return(nil)

				gh.EXPECT().Checkout(r0, "0").Return(nil, ErrCommitNotFound)

				gh.EXPECT().Checkout(r0, "1").Return(&v1.Commit{
					Sha: "sha1",
				}, nil)
				gt.EXPECT().Cover(r0).Return(&v1.Cover{
					Total: 0.2,
				}, nil)

				gh.EXPECT().Checkout(r1, "0").Return(&v1.Commit{
					Sha: "sha2",
				}, nil)
				gt.EXPECT().Cover(r1).Return(&v1.Cover{
					Total: 0.3,
				}, nil)

				gh.EXPECT().Checkout(r1, "1").Return(&v1.Commit{
					Sha: "sha3",
				}, nil)
				gt.EXPECT().Cover(r1).Return(&v1.Cover{
					Total: 0.4,
				}, nil)
			},
			want: &v1.Covalyzer{
				Repositories: []*v1.RepositoryCoverages{
					{
						Repository: r0,
						Coverages: []*v1.Coverage{
							{},
							{
								Commit: &v1.Commit{
									Sha: "sha1",
								},
								Cover: &v1.Cover{
									Total: 0.2,
								},
							},
						},
					},
					{
						Repository: r1,
						Coverages: []*v1.Coverage{
							{
								Commit: &v1.Commit{
									Sha: "sha2",
								},
								Cover: &v1.Cover{
									Total: 0.3,
								},
							},
							{
								Commit: &v1.Commit{
									Sha: "sha3",
								},
								Cover: &v1.Cover{
									Total: 0.4,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			gh := NewMockgh(ctrl)
			gt := NewMockgotool(ctrl)
			if tt.mock != nil {
				tt.mock(gh, gt)
			}
			c := NewCovalyzer(config, gh, gt)
			got, err := c.Run()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Errorf("Run() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

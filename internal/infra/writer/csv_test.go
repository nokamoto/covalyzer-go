package writer

import (
	"bytes"
	"encoding/csv"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

func TestCSVWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
		config  *v1.Config
		data    *v1.Covalyzer
		want    [][]string
		wantErr bool
	}{
		{
			name: "ok",
			config: &v1.Config{
				Timestamps: []string{"2024-01-01T00:00:00Z", "2024-02-01T00:00:00Z"},
			},
			data: &v1.Covalyzer{
				Repositories: []*v1.RepositoryCoverages{
					{
						Repository: &v1.Repository{
							Owner: "foo",
							Repo:  "bar",
						},
						Coverages: []*v1.Coverage{
							{
								Cover: &v1.Cover{
									Total: 0.5,
								},
							},
							{
								Cover: &v1.Cover{
									Total: 0.6,
								},
							},
						},
					},
					{
						Repository: &v1.Repository{
							Gh:    "example.com",
							Owner: "baz",
							Repo:  "qux",
						},
						Coverages: []*v1.Coverage{
							{
								Cover: &v1.Cover{
									Total: 0.7,
								},
							},
							{
								Cover: &v1.Cover{
									Total: 0.8,
								},
							},
						},
					},
				},
			},
			want: [][]string{
				{"github", "repository", "2024-01-01", "2024-02-01"},
				{"", "foo/bar", "0.5", "0.6"},
				{"example.com", "baz/qux", "0.7", "0.8"},
			},
		},
		{
			name: "failed to parse timestamp",
			config: &v1.Config{
				Timestamps: []string{"2024-01-01"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := &CSVWriter{
				file: t.TempDir() + "/covalyzer.csv",
			}
			err := sut.Write(tt.config, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}

			bs, _ := os.ReadFile(sut.file)
			r := csv.NewReader(bytes.NewBuffer(bs))
			got, _ := r.ReadAll()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

package config

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	"google.golang.org/protobuf/testing/protocmp"
)

//go:embed testdata/config.yaml
var configYAML []byte

func Test_NewConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  []byte
		want    *v1.Config
		wantErr bool
	}{
		{
			name:   "ok",
			config: configYAML,
			want: &v1.Config{
				Repositories: []*v1.Repository{
					{
						Owner: "nokamoto",
						Repo:  "covalyzer-go",
						GinkgoPackages: []string{
							"cmd/covalyzer-go-test",
						},
					},
				},
				Timestamps: []string{
					"2024-02-01T00:00:00Z",
					"2024-03-01T00:00:00Z",
				},
			},
		},
		{
			name:    "error if read invalid yaml",
			config:  []byte("invalid"),
			wantErr: true,
		},
		{
			name:    "error if protojson unmarshal fails",
			config:  []byte(`timestamps: "invalid type"`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempfile := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(tempfile, tt.config, 0644); err != nil {
				t.Fatal(err)
			}
			got, err := NewConfig(tempfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("NewConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

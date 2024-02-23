package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	"google.golang.org/protobuf/testing/protocmp"
)

func Test_NewConfig(t *testing.T) {
	testdata := `
repositories:
- owner: foo
  repo: bar
timestamps:
- 2024-01-01T00:00:00Z
`
	tempfile := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(tempfile, []byte(testdata), 0644); err != nil {
		t.Fatal(err)
	}

	expected := &v1.Config{
		Repositories: []*v1.Repository{
			{
				Owner: "foo",
				Repo:  "bar",
			},
		},
		Timestamps: []string{
			"2024-01-01T00:00:00Z",
		},
	}

	actual, err := NewConfig(tempfile)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		t.Errorf("NewConfig() mismatch (-want +got):\n%s", diff)
	}
}

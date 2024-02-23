package xslices

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Concat(t *testing.T) {
	tests := []struct {
		name string
		got  []string
		want []string
	}{
		{
			name: "empty",
			got:  Concat(),
		},
		{
			name: "string+string",
			got:  Concat("a", "b"),
			want: []string{"a", "b"},
		},
		{
			name: "slices+string+slices",
			got:  Concat([]string{"a"}, "b", []string{"c"}),
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tt.got); diff != "" {
				t.Errorf("Concat() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

package model

import "testing"

func TestChunkIDs(t *testing.T) {
	tests := []struct {
		name string
		ids  []string
		size int
		want [][]string
	}{
		{"empty", []string{}, 2, [][]string{}},
		{"exact", []string{"a", "b", "c", "d"}, 2, [][]string{{"a", "b"}, {"c", "d"}}},
		{"uneven", []string{"a", "b", "c"}, 2, [][]string{{"a", "b"}, {"c"}}},
		{"size zero", []string{"a", "b"}, 0, [][]string{{"a"}, {"b"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChunkIDs(tt.ids, tt.size)
			if len(got) != len(tt.want) {
				t.Fatalf("len=%d want %d", len(got), len(tt.want))
			}
			for i := range got {
				if len(got[i]) != len(tt.want[i]) {
					t.Fatalf("chunk %d len=%d want %d", i, len(got[i]), len(tt.want[i]))
				}
				for j := range got[i] {
					if got[i][j] != tt.want[i][j] {
						t.Fatalf("chunk %d[%d]=%q want %q", i, j, got[i][j], tt.want[i][j])
					}
				}
			}
		})
	}
}

func TestChunkIDsNoAliasing(t *testing.T) {
	ids := []string{"a", "b", "c", "d"}
	chunks := ChunkIDs(ids, 2)
	chunks[0][0] = "MUTATED"
	if ids[0] != "a" {
		t.Fatalf("mutating chunk corrupted input: ids[0]=%q", ids[0])
	}
}

package problems

import "testing"

func TestFindSmallestInt(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  int
	}{
		{
			name:  "simple case",
			slice: []int{3, 1, 4, 1, 5, 9, 2, 6},
			want:  1,
		},
		{
			name:  "with negative numbers",
			slice: []int{-1, -5, 0, 5, -10},
			want:  -10,
		},
		{
			name:  "empty slice",
			slice: []int{},
			want:  0, // Or handle as an error, for now, we expect 0 for an empty slice
		},
		{
			name:  "slice with one element",
			slice: []int{42},
			want:  42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindSmallestInt(tt.slice); got != tt.want {
				t.Errorf("FindSmallestInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

package problems

import "testing"

func TestFindSmallestIntInMap(t *testing.T) {
	tests := []struct {
		name string
		m    map[int]int
		want int
	}{
		{
			name: "simple case",
			m:    map[int]int{1: 3, 2: 1, 3: 4, 4: 1, 5: 5},
			want: 1,
		},
		{
			name: "with negative numbers",
			m:    map[int]int{1: -1, 2: -5, 3: 0, 4: 5, 5: -10},
			want: -10,
		},
		{
			name: "empty map",
			m:    map[int]int{},
			want: 0, // Or handle as an error, for now, we expect 0
		},
		{
			name: "map with one element",
			m:    map[int]int{1: 42},
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindSmallestIntInMap(tt.m); got != tt.want {
				t.Errorf("FindSmallestIntInMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

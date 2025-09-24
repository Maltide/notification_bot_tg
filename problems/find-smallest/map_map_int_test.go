package problems

import "testing"

func TestFindSmallestIntInMapMap(t *testing.T) {
	tests := []struct {
		name string
		m    map[int]map[int]int
		want int
	}{
		{
			name: "simple case",
			m: map[int]map[int]int{
				1: {1: 10, 2: 5},
				2: {1: 12, 2: 1},
				3: {1: 8, 2: 9},
			},
			want: 1,
		},
		{
			name: "with negative numbers",
			m: map[int]map[int]int{
				1: {1: -1, 2: -5},
				2: {1: 0, 2: -10},
			},
			want: -10,
		},
		{
			name: "empty map",
			m:    map[int]map[int]int{},
			want: 0, // Or handle as an error
		},
		{
			name: "map with one element",
			m:    map[int]map[int]int{1: {1: 42}},
			want: 42,
		},
		{
			name: "map with empty inner map",
			m: map[int]map[int]int{
				1: {},
				2: {1: 100},
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindSmallestIntInMapMap(tt.m); got != tt.want {
				t.Errorf("FindSmallestIntInMapMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

package problems

import (
	"testing"
	"time"
)

func TestFindSmallestTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		slice []time.Time
		want  time.Time
	}{
		{
			name: "simple case",
			slice: []time.Time{
				now.Add(3 * time.Hour),
				now.Add(1 * time.Hour),
				now.Add(4 * time.Hour),
			},
			want: now.Add(1 * time.Hour),
		},
		{
			name:  "empty slice",
			slice: []time.Time{},
			want:  time.Time{},
		},
		{
			name:  "slice with one element",
			slice: []time.Time{now},
			want:  now,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindSmallestTime(tt.slice); !got.Equal(tt.want) {
				t.Errorf("FindSmallestTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

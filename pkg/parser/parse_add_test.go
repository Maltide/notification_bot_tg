package parser

import (
	"testing"
	"time"
)

func TestParseAdd(t *testing.T) {
	now := time.Date(2025, 8, 31, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name, args, wantText string
		wantDue              time.Time
		wantErr              bool
	}{
		{
			name:     "ok simple",
			args:     "10m drink water",
			wantText: "drink water",
			wantDue:  now.Add(10 * time.Minute),
			wantErr:  false,
		},
		{
			name:    "err empty string",
			args:    "",
			wantErr: true,
		},
		{
			name:    "err only duration",
			args:    "10m",
			wantErr: true,
		},
		{
			name:    "err bad duration",
			args:    "abc task",
			wantErr: true,
		},
		{
			name:    "err zero duration",
			args:    "0m task",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			due, text, err := ParseAdd(tt.args, now)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if text != tt.wantText {
				t.Errorf("got text %q, want %q", text, tt.wantText)
			}
			if !due.Equal(tt.wantDue) {
				t.Errorf("got due %v, want %v", due, tt.wantDue)
			}
		})
	}
}

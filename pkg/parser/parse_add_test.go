package parser

import (
	"testing"
	"time"
)

func TestParseAddTask(t *testing.T) {

	parser := NewTimeParser()

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("Europe/Moscow", 3*60*60)
	}

	parser.NowFunc = func() time.Time {
		return time.Date(2025, 9, 29, 11, 10, 0, 0, loc)
	}

	tests := []struct {
		name, wantText string
		userinput      []string
		wantDue        time.Time
		wantErr        bool
	}{
		{
			name:      "ok simple",
			userinput: []string{"29.09.25", "11:10", "drink", "water"},
			wantText:  "drink water",
			wantDue:   time.Date(2025, 9, 29, 11, 10, 0, 0, loc),
			wantErr:   false,
		},
		{
			name:      "err empty string",
			userinput: []string{},
			wantErr:   true,
		},
		{
			name:      "err only duration",
			userinput: []string{"29.09.25", "11:00"},
			wantErr:   true,
		},
		{
			name:      "err bad duration",
			userinput: []string{"abc", "task", "time"},
			wantErr:   true,
		},
		{
			name:      "err zero duration",
			userinput: []string{"00.00.00", "00:00", "task"},
			wantErr:   true,
		},
		{
			name:      "day_check",
			userinput: []string{"30.09.25", "11:10", "daycheck", "test"},
			wantText:  "daycheck test",
			wantDue:   time.Date(2025, 9, 30, 11, 10, 0, 0, loc),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			due, text, err := parser.ParseAddTask(tt.userinput)

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

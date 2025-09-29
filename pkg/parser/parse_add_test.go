package parser

import (
	"testing"
	"time"
)

func TestParseAdd(t *testing.T) {

	now := time.Date(2025, 9, 29, 11, 00, 0, 0, Moscow_current_time)

	parser := NewTimeParser()

	tests := []struct {
		name, userinput, wantText string
		wantDue                   time.Time
		wantErr                   bool
	}{
		{
			name:      "ok simple",
			userinput: "29.09.25 11:10 drink water",
			wantText:  "drink water",
			wantDue:   time.Date(2025, 9, 29, 11, 10, 0, 0, Moscow_current_time),
			wantErr:   false,
		},
		{
			name:      "err empty string",
			userinput: "",
			wantErr:   true,
		},
		{
			name:      "err only duration",
			userinput: "29.09.25 11:00",
			wantErr:   true,
		},
		{
			name:      "err bad duration",
			userinput: "abc task",
			wantErr:   true,
		},
		{
			name:      "err zero duration",
			userinput: "00.00.00 00:00 task",
			wantErr:   true,
		},
		{
			name:      "day_check",
			userinput: "30.09.25 11:10 daycheck test",
			wantText:  "daycheck test",
			wantDue:   time.Date(2025, 9, 30, 11, 10, 0, 0, Moscow_current_time),
			wantErr:   false,
		},
		{
			name:      "1m",
			userinput: "29.09.25 11:01 go home",
			wantText:  "go home",
			wantDue:   time.Date(2025, 9, 29, 11, 01, 0, 0, Moscow_current_time),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			due, text, err := parser.ParseAddTask(tt.userinput, now)

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

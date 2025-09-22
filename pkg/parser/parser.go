package parser

import (
	"fmt"
	"strings"
	"time"
)

func NewTimeParser() *TimeParser {
	return &TimeParser{}
}

func (tp *TimeParser) ParseAddTask(args string, now time.Time) (due time.Time, text string, err error) {
	if strings.TrimSpace(args) == "" {
		return time.Time{}, "", fmt.Errorf("no task yet")
	}

	i := strings.Index(args, " ")
	if i == -1 {
		return time.Time{}, "", fmt.Errorf("duration and text are not created")
	}

	durationStr := args[:i]
	text = strings.TrimSpace(args[i+1:])

	dur, err := time.ParseDuration(durationStr)
	if err != nil || dur <= 0 {
		return time.Time{}, "", fmt.Errorf("duration time is zero or below")
	}
	due = now.Add(dur)
	return due, text, nil
}

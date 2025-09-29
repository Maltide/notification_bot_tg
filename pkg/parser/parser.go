package parser

import (
	"fmt"
	"strings"
	"time"
)

func NewTimeParser() *TimeParser {
	return &TimeParser{}
}

func (tp *TimeParser) ParseAddTask(userinput string) (due time.Time, text string, err error) {
	if strings.TrimSpace(userinput) == "" {
		return time.Time{}, "", fmt.Errorf("Заметка отсутствует")
	}

	shards := strings.SplitN(userinput, " ", 3)
	if len(shards) < 3 {
		return time.Time{}, "", fmt.Errorf("Нужен формат именно такой: DD.MM.YY HH:MM <text>")
	}

	notif_time := shards[0] + " " + shards[1]

	due, err = time.ParseInLocation(userinput, notif_time, time.Local)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("Неправильный формат даты или времени")
	}

	text = shards[2]

	return due, text, nil
}

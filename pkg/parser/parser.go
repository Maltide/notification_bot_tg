package parser

import (
	"fmt"
	"strings"
	"time"
)

var Moscow_current_time, _ = time.LoadLocation("Europe/Moscow")

const inputlayout = "02.01.06 15:04"

func NewTimeParser() *TimeParser {
	return &TimeParser{}
}

func (tp *TimeParser) ParseAddTask(userinput string) (due time.Time, text string, err error) {
	if strings.TrimSpace(userinput) == "" {
		return time.Time{}, "", fmt.Errorf("заметка отсутствует")
	}

	shards := strings.Fields(userinput)
	if len(shards) < 3 {
		return time.Time{}, "", fmt.Errorf("нужен формат именно такой: 12.11.25 15:05 Сходить в магазин")
	}

	notif_time := strings.ReplaceAll(shards[0]+" "+shards[1], ",", "")
	due, err = time.ParseInLocation(inputlayout, notif_time, Moscow_current_time)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("неправильный формат даты или времени")
	}

	if due.Before(time.Now()) {
		return time.Time{}, "", fmt.Errorf("пытаетесь ввести прошедшее время")
	}

	text = strings.Join(shards[:2], " ")

	return due, text, nil
}

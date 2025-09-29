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

func (tp *TimeParser) ParseAddTask(userinput string, now time.Time) (due time.Time, text string, err error) {
	if strings.TrimSpace(userinput) == "" {
		return time.Time{}, "", fmt.Errorf("Заметка отсутствует")
	}

	shards := strings.SplitN(userinput, " ", 3)
	if len(shards) < 3 {
		return time.Time{}, "", fmt.Errorf("Нужен формат именно такой: 12.11.25 15:05 <Сходить в магазин>")
	}

	notif_time := strings.ReplaceAll(shards[0]+" "+shards[1], ",", "")

	due, err = time.ParseInLocation(inputlayout, notif_time, Moscow_current_time)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("Неправильный формат даты или времени")
	}

	if due.Before(now) {
		return time.Time{}, "", fmt.Errorf("Вы пытаетесь ввести прошедшее время")
	}

	text = shards[2]

	return due, text, nil
}

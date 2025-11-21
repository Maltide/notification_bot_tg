package parser

import (
	"fmt"
	"strings"
	"time"
)

var moscow_current_time, _ = time.LoadLocation("Europe/Moscow")

const dateTimelayout = "02.01.06 15:04"

func NewTimeParser() *TimeParser {
	return &TimeParser{}
}

func (tp *TimeParser) ParseAddTask(userargs []string) (due time.Time, text string, err error) {

	if len(userargs) < 3 {
		return time.Time{}, "", fmt.Errorf("нужен формат именно такой: 12.11.25 15:05 Сходить в магазин")
	}

	notif_time := strings.ReplaceAll(userargs[0]+" "+userargs[1], ",", "")
	due, err = time.ParseInLocation(dateTimelayout, notif_time, moscow_current_time)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("неправильный формат даты или времени")
	}

	if due.Before(time.Now()) {
		return time.Time{}, "", fmt.Errorf("попытка ввести прошедшее время")
	}

	text = strings.Join(userargs[2:], " ")

	return due, text, nil
}

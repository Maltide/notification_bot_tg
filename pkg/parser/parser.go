package parser

import (
	"fmt"
	"strings"
	"time"
)

var moscowLocation = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		// In minimal Docker images (e.g. Alpine without tzdata) LoadLocation can fail.
		// Moscow doesn't use DST, so a fixed +03:00 offset is sufficient.
		return time.FixedZone("Europe/Moscow", 3*60*60)
	}
	return loc
}()

const dateTimelayout = "02.01.06 15:04"

// NewTimeParser constructs a parser that converts user input into due time and task text.
func NewTimeParser() *TimeParser {
	return &TimeParser{
		NowFunc: time.Now,
	}
}

// ParseAddTask parses arguments of the /add command and returns (due time, task text).
func (tp *TimeParser) ParseAddTask(userargs []string) (due time.Time, text string, err error) {
	if len(userargs) < 3 {
		return time.Time{}, "", fmt.Errorf("Ошибка формата\nПопробуйте использовать формат: /add 12.11.25 15:05 Сходить в магазин")
	}

	notifTime := strings.ReplaceAll(userargs[0]+" "+userargs[1], ",", "")
	due, err = time.ParseInLocation(dateTimelayout, notifTime, moscowLocation)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("неправильный формат даты или времени, посмотрите в /help")
	}

	if due.Before(tp.NowFunc()) {
		return time.Time{}, "", fmt.Errorf("попытка ввести прошедшее время")
	}

	text = strings.Join(userargs[2:], " ")

	return due, text, nil
}

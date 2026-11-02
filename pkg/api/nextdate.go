package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	dayRule   = "d"
	weekRule  = "w"
	monthRule = "m"
	yearRule  = "y"
)

var ErrRepeatRule = errors.New("incorrect the repeat rule")

var daysOfWeek = map[int]time.Weekday{
	1: time.Monday,
	2: time.Tuesday,
	3: time.Wednesday,
	4: time.Thursday,
	5: time.Friday,
	6: time.Saturday,
	7: time.Sunday,
}

var maxDaysInMonth = [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func convertToInt(s []string) ([]int, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("slice strings for converting to int is empty")
	}
	sliceInt := make([]int, len(s))
	for i, str := range s {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to convert string %q to int, got: %v", str, err)
		}
		sliceInt[i] = num
	}
	return sliceInt, nil
}

func firstDayOfMonth(date time.Time, day []int) bool {
	for _, num := range day {
		switch num {
		case -1:
			if date.AddDate(0, 0, 1).Day() == 1 {
				return true
			}
		case -2:
			if date.AddDate(0, 0, 2).Day() == 1 {
				return true
			}
		default:
			if date.Day() == num {
				return true
			}
		}
	}
	return false
}

func hasPossibleDayOfMonth(days, months []int) bool {
	for _, month := range months {
		for _, day := range days {
			if day <= maxDaysInMonth[month - 1] {
				return true
			}
		}
	}
	return false
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf(`date "%s" is not a valid date format`, dstart)
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty %q", repeat)
	}

	interval := strings.Split(repeat, " ")

	var nextDate string
	switch interval[0] {
	case dayRule:
		if len(interval) != 2 {
			return "", ErrRepeatRule
		}

		day, err := strconv.Atoi(interval[1])
		if err != nil {
			return "", fmt.Errorf("failed to convert day %q to int, got: %v", interval[1], err)
		}

		if day <= 0 || day > 400 {
			return "", fmt.Errorf("day must be positive or greater than 0 and not must greater than 400, got %d", day)
		}

		date = date.AddDate(0, 0, day)
		for !date.After(now) {
			date = date.AddDate(0, 0, day)
		}
		nextDate = date.Format("20060102")

	case weekRule:
		if len(interval) != 2 {
			return "", ErrRepeatRule
		}

		weekDays := strings.Split(interval[1], ",")
		weekDaysInt, err := convertToInt(weekDays)
		if err != nil {
			return "", fmt.Errorf("failed to convert week days to int: %w", err)
		}

		// в тз указано, что задача назначается в указанные дни недели по отношению к now
		// берем первый указанный день из списка, например, если now пн и в переданном правиле указаны w 1,2,
		// то ближайший понедельник к понедельнику, который now, через неделю
		// в некорректном правиле, например указанные дни недели 3,8, берем первый указанный день, он корректный
		// вернется ближайшая среда от переданного в now дня, 8 просто опускается
		nextDay := weekDaysInt[0]

		_, ok := daysOfWeek[nextDay]
		if !ok {
			return "", fmt.Errorf("day doesn't match the week")
		}

		date = date.AddDate(0, 0, 1)
		for !date.After(now) || date.Weekday() != daysOfWeek[nextDay] {
			date = date.AddDate(0, 0, 1)
		}
		nextDate = date.Format("20060102")

	case monthRule:
		switch len(interval) {
		case 2:
			days := strings.Split(interval[1], ",")
			numDays, err := convertToInt(days)
			if err != nil {
				return "", fmt.Errorf("failed to convert days to int: %w", err)
			}

			for _, day := range numDays {
				if day < -2 || day == 0 || day > 31 {
					return "", fmt.Errorf("day out of range, min day number is %d", day)
				}
			}

			for !date.After(now) || !firstDayOfMonth(date, numDays) {
				date = date.AddDate(0, 0, 1)
			}
			nextDate = date.Format("20060102")

		case 3:
			days := strings.Split(interval[1], ",")
			month := strings.Split(interval[2], ",")

			numDays, err := convertToInt(days)
			if err != nil {
				return "", fmt.Errorf("failed to convert days to int: %w", err)
			}

			for _, day := range numDays {
				if day < -2 || day == 0 || day > 31 {
					return "", fmt.Errorf("day out of range, min day number is %d", day)
				}
			}

			months, err := convertToInt(month)
			if err != nil {
				return "", fmt.Errorf("failed to convert month to int: %w", err)
			}

			allowed := make(map[time.Month]bool)
			for _, m := range months {
				if m < 1 || m > 12 {
					return "", fmt.Errorf("month must be between 1 and 12, got %d", m)
				}
				allowed[time.Month(m)] = true
			}

			if !hasPossibleDayOfMonth(numDays, months) {
				return "", fmt.Errorf("no possible day %v in months %v", numDays, months)
			}

			for !date.After(now) || !firstDayOfMonth(date, numDays) || !allowed[date.Month()] {
				date = date.AddDate(0, 0, 1)
			}
			nextDate = date.Format("20060102")

		default:
			return "", ErrRepeatRule
		}

	case yearRule:
		date = date.AddDate(1, 0, 0)
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
		}
		nextDate = date.Format("20060102")
	default:
		return "", ErrRepeatRule
	}
	return nextDate, nil
}

func getNextDate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	now := query.Get("now")
	if now == "" {
		now = time.Now().Format("20060102")
	}

	dateStart := query.Get("date")
	if dateStart == "" {
		http.Error(w, "date is empty", http.StatusBadRequest)
		return
	}

	repeat := query.Get("repeat")
	if repeat == "" {
		http.Error(w, "repeat is empty", http.StatusBadRequest)
		return
	}

	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, "failed to parse time now", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(nowTime, dateStart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(nextDate)); err != nil {
		log.Printf("failed to write response body for client: %v ", err)
	}
}

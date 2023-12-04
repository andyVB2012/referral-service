package util

import "time"

func DateFirstAndLastMonthMoment(year int, month int) (time.Time, time.Time) {
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1).Add(time.Second * (82800 + 3540 + 59))
	return firstOfMonth, lastOfMonth
}

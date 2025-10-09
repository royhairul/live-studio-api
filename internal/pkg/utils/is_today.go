package utils

import "time"

func IsToday(datetime int64) bool {
	startTime := time.UnixMilli(datetime)
	now := time.Now()

	y1, m1, d1 := startTime.Date()
	y2, m2, d2 := now.Date()

	return (y1 == y2 && m1 == m2 && d1 == d2)
}

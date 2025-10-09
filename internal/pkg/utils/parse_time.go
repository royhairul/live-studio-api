package utils

import (
	"fmt"
	"time"
)

func ParseTime(timeStr string) (time.Time, error) {
	layout := "15:05"
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{}, fmt.Errorf("error load location time: %w", err)
	}

	parsed, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("error format time: %w", err)
	}

	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, loc), nil
}

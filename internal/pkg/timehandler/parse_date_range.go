package timehandler

import (
	"fmt"
	"time"
)

// ParseDateRange received startDate & endDate string (format: yyyy-mm-dd)
// return start with 00:00:00 until end to 23:59:59
func ParseDateRange(startDate string, endDate string) (*time.Time, *time.Time, error) {
	start, err := ParseDate(startDate)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid start date: %w", err)
	}
	start = StartOfDay(*start)

	end, err := ParseDate(endDate)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid end date: %w", err)
	}
	end = EndOfDay(*end)

	return start, end, nil
}

// StartOfDay set time to 00:00:00
func StartOfDay(t time.Time) *time.Time {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &start
}

// EndOfDay set time to 23:59:59.999999999
func EndOfDay(t time.Time) *time.Time {
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
	return &end
}

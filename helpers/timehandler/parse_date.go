package timehandler

import (
	"fmt"
	"time"
)

func ParseDate(str string) (*time.Time, error) {
	layout := "2006-01-02"

	parsedDate, err := time.Parse(layout, str)
	if err != nil {
		return nil, fmt.Errorf("Error parsing date: %w", err)
	}

	return &parsedDate, nil
}

func ParseInt64Date(timestamp int64) *time.Time {
	if timestamp <= 0 {
		return nil
	}

	t := time.Unix(timestamp, 0).UTC()
	return &t
}

func FormatDate(date *time.Time) string {
	layout := "2006-01-02"

	return date.Format(layout)
}

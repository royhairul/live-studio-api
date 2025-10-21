package timehandler

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/internal/pkg/constants"
)

// ParseDate parses a date string into a time.Time object based on the defined layout.
func ParseDate(str string) (*time.Time, error) {
	layout := constants.LayoutYYMMDD

	parsedDate, err := time.Parse(layout, str)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	return &parsedDate, nil
}

// ParseInt64Date converts a Unix timestamp (in seconds) into a *time.Time object (in UTC).
func ParseInt64Date(timestamp int64) *time.Time {
	if timestamp <= 0 {
		return nil
	}

	t := time.Unix(timestamp, 0).UTC()
	return &t
}

// FormatDate formats a *time.Time object into a string using the defined layout.
func FormatDate(date *time.Time) string {
	if date == nil {
		return ""
	}

	layout := constants.LayoutYYMMDD
	return date.Format(layout)
}

// FormatInt64Date converts a Unix timestamp (in seconds) directly into a formatted date string.
func FormatInt64Date(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}

	layout := constants.LayoutYYMMDD
	return time.Unix(timestamp, 0).UTC().Format(layout)
}

package timehandler

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/internal/pkg/constants"
)

func ParseDate(str string) (*time.Time, error) {
	layout := constants.LayoutYYMMDD

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
	layout := constants.LayoutYYMMDD

	return date.Format(layout)
}

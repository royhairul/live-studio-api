package timehandler

import (
	"fmt"
	"time"
)

func ParseDateRange(startDate string, endDate string) (*time.Time, *time.Time, error) {
	start, err := ParseDate(startDate)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid start date: %w", err)
	}

	end, err := ParseDate(endDate)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid end date: %w", err)
	}

	return start, end, nil
}

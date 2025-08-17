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

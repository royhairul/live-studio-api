package params

import "time"

type AccountadsFilter struct {
	ID        *string
	AccountID *string
	StartDate *time.Time
	EndDate   *time.Time
}

package params

import "time"

type AccountadsFilter struct {
	ID         *string
	AccountID  *string
	AccountIDs []string
	StartDate  *time.Time
	EndDate    *time.Time
}

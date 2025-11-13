package params

import "time"

type TransactionFilter struct {
	ID        *string
	UniqueID  *string
	Status    *string
	AccountID *string
	StudioID  *string
	StartTime *time.Time
	EndTime   *time.Time
}

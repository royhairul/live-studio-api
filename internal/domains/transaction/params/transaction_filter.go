package params

import "time"

type TransactionFilter struct {
	ID        *string
	Status    *string
	AccountID *string
	StartTime *time.Time
	EndTime   *time.Time
}

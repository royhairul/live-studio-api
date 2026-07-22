package params

import "time"

type LiveFilter struct {
	ID        *string
	SessionID *int64
	AccountID *string
	StudioID  *string
	StartTime *time.Time
	EndTime   *time.Time

	// Pagination. PageSize <= 0 means "no limit".
	Page     int
	PageSize int
}

func (f LiveFilter) Offset() int {
	if f.Page <= 1 {
		return 0
	}
	return (f.Page - 1) * f.PageSize
}

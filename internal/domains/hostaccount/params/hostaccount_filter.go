package params

import "time"

type HostAccountFilter struct {
	ID        *string
	HostID    *string
	AccountID *string
	StudioID  *string

	// ActiveAt keeps only assignments covering this instant.
	ActiveAt *time.Time

	// StartDate/EndDate keep assignments overlapping the window, which is what a
	// report over a date range needs — an assignment that ended mid-range still
	// owns the sessions it covered.
	StartDate *time.Time
	EndDate   *time.Time

	// OnlyActive keeps only assignments that have not been ended (valid_to IS NULL).
	OnlyActive bool

	// ExcludeID drops one row, so an update can check for conflicts without
	// matching itself.
	ExcludeID *uint
}

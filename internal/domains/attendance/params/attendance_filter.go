package params

import "time"

type AttendanceFilter struct {
	AccountID *string
	HostID    *string
	ShiftID   *string
	StudioID  *string
	StartTime *time.Time
	EndTime   *time.Time
}

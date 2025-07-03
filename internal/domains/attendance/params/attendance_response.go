package params

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceResponse struct {
	ID uint

	HostID uuid.UUID `json:"host_id"`
	Name   string    `json:"host_name"`

	Date     *time.Time `json:"date"`
	CheckIn  *time.Time `json:"check_in"`
	CheckOut *time.Time `json:"check_out"`

	ShiftStartTime time.Time `json:"start_time"`
	ShiftEndTime   time.Time `json:"end_time"`

	Note string `json:"note"`
}

type AttendanceCheckInResult struct {
	HostName string `json:"host_name"`
	Message  string `json:"message"`
	Status   string `json:"status"` // success | failed | warning
}

type AttendanceCheckInSummary struct {
	Message      string                    `json:"message"` // contoh: "Berhasil check-in semua host" atau "Sebagian berhasil check-in"
	SuccessCount int                       `json:"success_count"`
	FailedCount  int                       `json:"failed_count"`
	Results      []AttendanceCheckInResult `json:"results"`
}

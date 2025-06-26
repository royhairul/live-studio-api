package dto

type CheckAttendanceRequestDTO struct {
	HostID  uint   `json:"host_id" binding:"required"`
	Date    string `json:"date" binding:"required"` // ISO 8601 string format (ex: 2025-06-14T00:00:00Z)
	ShiftID uint   `json:"shift_id" binding:"required"`
}

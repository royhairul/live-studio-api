package params

type AttendanceCheckInRequest struct {
	Date     string `json:"date" vaidate:"required"`
	HostID   string `json:"host_id" validate:"required"`
	ShiftID  uint   `json:"shift_id" validate:"required"`
	StudioID uint   `json:"studio_id" validate:"required"`
}

type AttendanceCheckOutRequest struct {
	ID uint `json:"id" validate:"required"`
}

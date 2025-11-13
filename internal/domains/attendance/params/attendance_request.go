package params

type AttendanceCheckInRequest struct {
	HostID   string `json:"host_id" validate:"required"`
	ShiftID  uint   `json:"shift_id,string" validate:"required"`
	StudioID uint   `json:"studio_id,string" validate:"required"`
}

type AttendanceCheckOutRequest struct {
	ID []uint `json:"id" validate:"required"`
}

package dto

type CreateHostScheduleDTO struct {
	HostID  uint   `json:"host_id" binding:"required"`
	ShiftID uint   `json:"shift_id" binding:"required"`
	Date    string `json:"date" binding:"required"`
}

type UpdateHostScheduleDTO struct {
	HostID  *uint   `json:"host_id"`
	ShiftID *uint   `json:"shift_id"`
	Date    *string `json:"date"`
}

type SwitchHostScheduleDTO struct {
	FromHostID     uint `json:"from_host_id" binding:"required"`
	FromScheduleID uint `json:"from_schedule_id" binding:"required"`
	ToHostID       uint `json:"to_host_id" binding:"required"`
	ToScheduleID   uint `json:"to_schedule_id" binding:"required"`
}
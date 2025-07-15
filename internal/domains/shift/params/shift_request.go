package params

type ShiftRequest struct {
	// TODO: add request fields
}

type CreateShiftRequest struct {
	Name      string `json:"name" validate:"required"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

type UpdateShiftRequest struct {
	Name      *string `json:"name" validate:"omitempty"`
	StartTime *string `json:"start_time" validate:"omitempty"`
	EndTime   *string `json:"end_time" validate:"omitempty"`
}

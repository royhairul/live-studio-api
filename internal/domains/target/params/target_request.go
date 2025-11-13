package params

type TargetRequest struct {
	Month  string `json:"month"`
	Year   string `json:"year"`
	Studio string `json:"studio"`
}

type CreateTargetRequest struct {
	StudioID     uint   `json:"studio_id" validate:"required"`
	Date         string `json:"date" validate:"required"`
	TargetGMV    int64  `json:"target_gmv" validate:"required"`
	TargetIncome int64  `json:"target_income" validate:"required"`
}

type UpdateTargetRequest struct {
	// TODO: add request fields
	StudioID     uint    `json:"studio_id"`
	Date         *string `json:"date" validate:"omitempty"`
	TargetGMV    *int64  `json:"target_gmv" validate:"omitempty"`
	TargetIncome *int64  `json:"target_income" validate:"omitempty"`
}

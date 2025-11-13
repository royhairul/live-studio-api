package params

type AccountInfoRequest struct {
	Cookie string `json:"cookie"`
}

type CreateAccountRequest struct {
	StudioID uint16 `json:"studio_id" validate:"required"`
	Cookie   string `json:"cookie" validate:"required"`
	Device   string `json:"device" validate:"omitempty"`
}

type UpdateAccountRequest struct {
	StudioID *uint16 `json:"studio_id" validate:"omitempty,min=1"`
	Cookie   *string `json:"cookie" validate:"omitempty,min=1"`
	Device   string  `json:"device" validate:"omitempty"`
}

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

// PatchAccountRequest is used for PATCH /account/:id - partial updates
// All fields are optional pointers to allow partial updates
type PatchAccountRequest struct {
	IsActive *bool   `json:"is_active,omitempty"`
	StudioID *uint16 `json:"studio_id,omitempty" validate:"omitempty,min=1"`
	Cookie   *string `json:"cookie,omitempty" validate:"omitempty,min=1"`
	Device   *string `json:"device,omitempty"`
}

package params

type CreateHostRequest struct {
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	StudioID uint   `json:"studio_id" validate:"required"`
}

type UpdateHostRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=1"`
	Phone    *string `json:"phone" validate:"omitempty,min=1"`
	StudioID *uint   `json:"studio_id" validate:"omitempty,min=1"`
}

package params

type StudioRequest struct {
	// TODO: add request fields
}

type CreateStudioRequest struct {
	Name    string `json:"name" validate:"required,min=4"`
	Address string `json:"address" validate:"omitempty"`
}

type UpdateStudioRequest struct {
	Name    *string `json:"name" validate:"omitempty,min=4"`
	Address *string `json:"address" validate:"omitempty"`
}

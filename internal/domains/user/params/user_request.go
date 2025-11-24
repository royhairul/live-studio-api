package params

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email,min=3"`
	Password string `json:"password" validate:"required,min=3"`
	RoleID   uint   `json:"roleID" validate:"required"`

	Phone    string `json:"phone" validate:"omitempty,min=1"`
	StudioID uint   `json:"studio_id" validate:"omitempty,min=1"`
}

type UpdateUserRequest struct {
	Name   string `json:"name" validate:"required,min=3"`
	Email  string `json:"email" validate:"required,email,min=3"`
	RoleID *uint  `json:"roleID"`

	Phone    string `json:"phone" validate:"omitempty,min=1"`
	StudioID uint   `json:"studio_id" validate:"omitempty,min=1"`
}

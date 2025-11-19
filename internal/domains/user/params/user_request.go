package params

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email,min=3"`
	Password string `json:"password" validate:"required,min=3"`
	RoleID   uint   `json:"roleID" validate:"required"`
}

type UpdateUserRequest struct {
	Name   string `json:"name" validate:"required,min=3"`
	Email  string `json:"email" validate:"required,email,min=3"`
	RoleID *uint  `json:"roleID"`
}

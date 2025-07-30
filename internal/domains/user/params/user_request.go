package params

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email,min=3"`
	Password string `json:"password" binding:"required,min=3"`
	RoleID   uint   `json:"roleID" binding:"required"`
}

type UpdateUserRequest struct {
	Name   string `json:"name" binding:"required,min=3"`
	Email  string `json:"email" binding:"required,email,min=3"`
	RoleID *uint  `json:"roleID"`
}

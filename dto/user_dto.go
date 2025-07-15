package dto

type LoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterDTO struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   uint   `json:"role" binding:"required"`
}

type CreateUserDTO struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   uint   `json:"role" binding:"required"`
	Name     string `json:"name" binding:"required"`
	StudioID uint   `json:"studio_id"` // hanya untuk role host, studio_id harus diisi
	Phone    string `json:"phone"`     // hanya untuk role host, phone_number harus diisi
}

type UpdateUserDTO struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	RoleID      uint   `json:"role"`
	CurrentUser uint   `json:"current_user"` // username dari user yang membuat
}

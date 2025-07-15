package dto

type CreateRoleDTO struct {
	Name        string `json:"name" validate:"required"`
	Permissions []uint `json:"permissions" validate:"required"`
}

type UpdateRoleDTO struct {
	Name        *string `json:"name"`
	Permissions *[]uint `json:"permissions"`
}

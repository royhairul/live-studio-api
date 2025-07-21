package params

type RoleRequest struct {
	// TODO: add request fields
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Permissions []uint `json:"permissions" validate:"required"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name" validate:"omitempty"`
	Permissions *[]uint `json:"permissions" validate:"omitempty"`
}

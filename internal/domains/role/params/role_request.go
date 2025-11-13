package params

type RoleRequest struct {
	// TODO: add request fields
}

type CreateRoleRequest struct {
	// TODO: add request fields
	Name        string `json:"name" validate:"required"`
	Permissions []uint `json:"permissions" validate:"required"`
}

type UpdateRoleRequest struct {
	// TODO: add request fields
	Name        string `json:"name,omitempty"`
	Permissions []uint `json:"permissions,omitempty"`
}

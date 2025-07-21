package params

type PermissionRequest struct {
	// TODO: add request fields
}

type CreatePermissionRequest struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        *string `json:"name"`
	Group       *string `json:"group"`
	Description *string `json:"description"`
}

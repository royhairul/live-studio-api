package params

type PermissionRequest struct {
	// TODO: add request fields
}

type CreatePermissionRequest struct {
	// TODO: add request fields
	Name        string `json:"name" gorm:"unique"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	// TODO: add request fields
}

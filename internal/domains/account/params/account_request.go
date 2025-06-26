package params

type AccountInfoRequest struct {
	Cookie string `json:"cookie"`
}

type CreateAccountRequest struct {
	StudioID uint16 `json:"studio_id" binding:"required"`
	Cookie   string `json:"cookie" binding:"required"`
}

type UpdateAccountRequest struct {
	StudioID *uint16 `json:"studio_id"`
	Cookie   *string `json:"cookie"`
}

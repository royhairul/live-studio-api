package dto

type CreateHostDTO struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required" validate:"phoneid"`
	StudioID uint16 `json:"studio_id" binding:"required"`
}

type UpdateHostDTO struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone" validate:"phoneid"`
	StudioID *uint16 `json:"studio_id"`
}
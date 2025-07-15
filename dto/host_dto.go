package dto

type HostResponse struct {
	ID         uint   `json:"ID"`
	Name       string `json:"Name"`
	Phone      string `json:"Phone"`
	StudioID   uint16 `json:"studio_id"`
	StudioName string `json:"StudioName"`
}

type CreateHostDTO struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required" validate:"phoneid"`
	StudioID uint16 `json:"studio_id" binding:"required"`
	UserID   uint16 `json:"user_id" binding:"required"`
}

type UpdateHostDTO struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone" validate:"phoneid"`
	StudioID *uint16 `json:"studio_id"`
}

package dto

type Studio struct {
	Name string `json:"studioname" binding:"required"`
}
package dto

type Studio struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}
package dto

type ForgotPasswordDTO struct {
	Email string `json:"email" binding:"required"`
}
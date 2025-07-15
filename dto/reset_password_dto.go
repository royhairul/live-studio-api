package dto

type ResetPassword struct {
	Otp             string `json:"otp" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
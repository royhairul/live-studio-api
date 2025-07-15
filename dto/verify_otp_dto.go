package dto

type VerifyOtp struct {
	Otp string `json:"otp" binding:"required"`
}
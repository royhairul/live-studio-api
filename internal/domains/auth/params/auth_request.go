package params

type AuthBase struct {
	Email    string `json:"email" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=3"`
}

type LoginRequest struct {
	AuthBase
}

type RegisterRequest struct {
	AuthBase
	Name string `json:"name" binding:"required,min=3"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,min=3"`
}

type ResetPasswordRequest struct {
	OTP             string `json:"otp" binding:"required,min=3"`
	Password        string `json:"password" binding:"required,min=3"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=3"`
}

type VerifyOTPRequest struct {
	OTP string `json:"otp" binding:"required,min=3"`
}

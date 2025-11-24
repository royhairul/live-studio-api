package params

type AuthBase struct {
	Email    string `json:"email" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=3"`
}

type LoginRequest struct {
	AuthBase
}

type RegisterRequest struct {
	AuthBase
	Name string `json:"name" validate:"required,min=3"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,min=3"`
}

type ResetPasswordRequest struct {
	OTP             string `json:"otp" validate:"required,min=3"`
	Password        string `json:"password" validate:"required,min=3"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=3"`
}

type VerifyOTPRequest struct {
	OTP string `json:"otp" validate:"required,min=3"`
}

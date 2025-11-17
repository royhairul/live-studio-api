package params

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type RegisterResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ForgotPasswordResponse struct {
	Email string `json:"email"`
	Otp   string
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}

type MeResponse struct {
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	TenantID    string   `json:"tenant_id"`
	Permissions []string `json:"permissions"`
}

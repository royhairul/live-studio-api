package dto

type CreateAccountCookiesDTO struct {
	StudioID uint16 `json:"studio_id" binding:"required"`
	Cookies  string `json:"cookies" binding:"required"`
}

type CreateAccountDTO struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

package dto

type CreateAccountCookiesDTO struct {
	Cookies string `json:"cookies" binding:"required"`
}

type CreateAccountDTO struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
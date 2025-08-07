package params

import "github.com/royhairul/live-studio-api/internal/domains/account/entity"

type AccountResponse struct {
	ID         uint   `json:"id"`
	UniqueID   string `json:"unique_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Platform   string `json:"platform"`
	StudioName string `json:"studio_name"`
	Cookie     string `json:"cookie"`
}

func NewAccountResponse(account *entity.Account) *AccountResponse {
	return &AccountResponse{
		ID:         account.ID,
		UniqueID:   account.UniqueID,
		Name:       account.Name,
		Username:   account.Username,
		Email:      account.Email,
		Platform:   account.Platform,
		StudioName: account.Studio.Name,
		Cookie:     account.Cookie,
	}
}

type AccountDetailResponse struct {
	ID         uint   `json:"id"`
	UniqueID   string `json:"unique_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Platform   string `json:"platform"`
	StudioName string `json:"studio_name"`
	Cookie     string `json:"cookie"`
}

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
}

func NewAccountResponse(accounts []*entity.Account) []*AccountResponse {
	accountsResp := []*AccountResponse{}
	for _, acc := range accounts {
		accountsResp = append(accountsResp, &AccountResponse{
			ID:         acc.ID,
			UniqueID:   acc.UniqueID,
			Name:       acc.Name,
			Username:   acc.Username,
			Email:      acc.Email,
			Platform:   acc.Platform,
			StudioName: acc.Studio.Name,
		})
	}
	return accountsResp
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

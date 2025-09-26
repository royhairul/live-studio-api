package params

import (
	"github.com/royhairul/live-studio-api/internal/domains/accountads/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/constants"
)

type AccountadsResponse struct {
	// TODO: add response fields
	ID   uint   `json:"id"`
	Date string `json:"date"`
	Ads  uint   `json:"ads"`
	Helo string `json:"helo"`
}

func NewAccountadsResponse(accountAds *entity.Accountads) *AccountadsResponse {
	return &AccountadsResponse{
		ID:   accountAds.ID,
		Date: accountAds.Date.Format(constants.LayoutMMYY),
		Ads:  accountAds.Spend,
		Helo: "ahiii",
	}
}

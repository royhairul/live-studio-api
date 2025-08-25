package params

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/accountads/entity"
)

type AccountadsResponse struct {
	// TODO: add response fields
	ID   uint       `json:"id"`
	Date *time.Time `json:"date"`
	Ads  uint       `json:"ads"`
}

func NewAccountadsResponse(accountAds *entity.Accountads) *AccountadsResponse {
	return &AccountadsResponse{
		ID:   accountAds.ID,
		Date: accountAds.Date,
		Ads:  accountAds.Spend,
	}
}

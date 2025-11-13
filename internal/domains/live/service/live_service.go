package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/live/params"
)

type LiveService interface {
	GetLive() ([]*params.LiveResponse, error)
	GetLiveDetail(accountID, sessionID, productPage, productPageSize string) (*params.LiveDetailResponse, error)
}

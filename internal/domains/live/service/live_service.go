package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/live/params"
)

type LiveService interface {
	GetLive(ctx context.Context) ([]*params.LiveResponse, error)
	GetLiveDetail(ctx context.Context, accountID, sessionID, productPage, productPageSize string) (*params.LiveDetailResponse, error)
}

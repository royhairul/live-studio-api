package service

import (
	"context"

	shopeeparams "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"github.com/royhairul/live-studio-api/internal/domains/live/params"
)

type LiveService interface {
	// Preview: realtime, streamed over WebSocket, ongoing sessions only.
	GetLive(ctx context.Context) ([]*params.LiveResponse, error)
	GetLiveDetail(ctx context.Context, accountID, sessionID, productPage, productPageSize string) (*params.LiveDetailResponse, error)

	// History: served from the database, never proxied from Shopee.
	GetStoredHistory(ctx context.Context, filter params.LiveFilter) (*params.StoredLiveResponse, error)

	SyncHistory(ctx context.Context, accountID string, req shopeeparams.ShopeeLiveHistoryRequest) (*params.LiveSyncResponse, error)
	SyncAllHistory(ctx context.Context, req shopeeparams.ShopeeLiveHistoryRequest) ([]*params.LiveSyncResponse, error)
}

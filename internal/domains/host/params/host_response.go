package params

import (
	"github.com/google/uuid"
	"github.com/royhairul/live-studio-api/internal/domains/host/entity"
)

type HostResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	StudioID   uint      `json:"studio_id"`
	StudioName string    `json:"studio_name"`
}

func NewHostResponse(host *entity.Host) *HostResponse {
	return &HostResponse{
		ID:         *host.ID,
		Name:       host.Name,
		Phone:      host.Phone,
		StudioID:   host.Studio.ID,
		StudioName: host.Studio.Name,
	}
}

type HostGroupedByStudioResponse struct {
	StudioID   uint           `json:"studio_id"`
	StudioName string         `json:"studio_name"`
	Hosts      []HostResponse `json:"hosts"`
}

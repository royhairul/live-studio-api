package params

import "github.com/royhairul/live-studio-api/internal/domains/studio/entity"

type StudioResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func NewStudioResponse(studio *entity.Studio) *StudioResponse {
	return &StudioResponse{
		ID:      studio.ID,
		Name:    studio.Name,
		Address: studio.Address,
	}
}

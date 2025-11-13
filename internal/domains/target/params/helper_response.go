package params

import (
	"fmt"

	"github.com/royhairul/live-studio-api/internal/domains/target/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/constants"
)

func NewTargetResponse(t *entity.Target) *TargetResponse {
	return &TargetResponse{
		StudioID:   fmt.Sprintf("%d", t.StudioID),
		StudioName: t.Studio.Name,
		Date:       t.Date.Format(constants.LayoutYYMMDD),
	}
}

func NewCreatedTargetResponse(t *entity.Target) *CreatedTargetResponse {
	return &CreatedTargetResponse{
		StudioID:     fmt.Sprintf("%d", t.StudioID),
		StudioName:   t.Studio.Name,
		Date:         t.Date.Format(constants.LayoutYYMMDD),
		TargetGMV:    t.TargetGMV,
		TargetIncome: t.TargetIncome,
	}
}

func NewUpdatedTargetResponse(t *entity.Target) *UpdatedTargetResponse {
	return &UpdatedTargetResponse{
		StudioID:     fmt.Sprintf("%d", t.StudioID),
		StudioName:   t.Studio.Name,
		Date:         t.Date.Format(constants.LayoutYYMMDD),
		TargetGMV:    t.TargetGMV,
		TargetIncome: t.TargetIncome,
	}
}

package params

import (
	"fmt"

	"github.com/royhairul/live-studio-api/internal/domains/target/entity"
)

type Metric struct {
	Real   int64   `json:"real"`
	Target int64   `json:"target"`
	Ratio  float64 `json:"ratio"`
}

type CreatedTargetResponse struct {
	StudioID     string `json:"studio_id"`
	StudioName   string `json:"studio_name"`
	Date         string `json:"date"`
	TargetGMV    int64  `json:"target_gmv"`
	TargetIncome int64  `json:"target_income"`
}

func NewCreatedTargetResponse(t *entity.Target) *CreatedTargetResponse {
	return &CreatedTargetResponse{
		StudioID:     fmt.Sprintf("%d", t.StudioID),
		StudioName:   t.Studio.Name, // pastikan preload
		Date:         t.Date.Format("2006-01-02"),
		TargetGMV:    t.TargetGMV,
		TargetIncome: t.TargetIncome,
	}
}

type UpdatedTargetResponse struct {
	StudioID     string `json:"studio_id"`
	StudioName   string `json:"studio_name"`
	Date         string `json:"date"`
	TargetGMV    int64  `json:"target_gmv"`
	TargetIncome int64  `json:"target_income"`
}

func NewUpdatedTargetResponse(t *entity.Target) *UpdatedTargetResponse {
	return &UpdatedTargetResponse{
		StudioID:     fmt.Sprintf("%d", t.StudioID),
		StudioName:   t.Studio.Name, // pastikan preload
		Date:         t.Date.Format("2006-01-02"),
		TargetGMV:    t.TargetGMV,
		TargetIncome: t.TargetIncome,
	}
}

type TargetResponse struct {
	// TODO: add response fields
	StudioID   string `json:"studio_id"`
	StudioName string `json:"studio_name"`
	Date       string `json:"date"`
	GMV        Metric `json:"gmv"`
	Income     Metric `json:"income"`
}

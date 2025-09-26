package params

import (
	performaparams "github.com/royhairul/live-studio-api/internal/domains/performa/params"
)

type PeriodInfo struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Days  int    `json:"days"`
}

type Metric struct {
	Total int64 `json:"total"`
	Diff  int64 `json:"diff"`
	Ratio int64 `json:"ratio"`
}

type Metrics struct {
	Commission Metric `json:"commission"`
	Income     Metric `json:"income"`
	GMV        Metric `json:"gmv"`
	Ads        Metric `json:"ads"`
	Account    int64  `json:"account"`
	Host       int64  `json:"host"`
}

type DashboardResponse struct {
	CurrentPeriod  PeriodInfo                                  `json:"current_period"`
	PreviousPeriod PeriodInfo                                  `json:"previous_period"`
	Metrics        Metrics                                     `json:"metrics"`
	List           []performaparams.PerformaStudioItemResponse `json:"list"`
}

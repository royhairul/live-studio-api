package params

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
}

type PerformaStudioResponse struct {
	CurrentPeriod  PeriodInfo                   `json:"current_period"`
	PreviousPeriod PeriodInfo                   `json:"previous_period"`
	Metrics        Metrics                      `json:"metrics"`
	List           []PerformaStudioItemResponse `json:"list"`
}

type PerformaStudioItemResponse struct {
	StudioID   string `json:"studio_id"`
	StudioName string `json:"studio_name"`
	Income     int64  `json:"income"`
	Commission int64  `json:"commission"`
	GMV        int64  `json:"gmv"`
	Ads        int64  `json:"ads"`
}

type PerformaStudioDetailResponse struct {
	StudioID       uint                               `json:"studio_id"`
	StudioName     string                             `json:"studio_name"`
	CurrentPeriod  PeriodInfo                         `json:"current_period"`
	PreviousPeriod PeriodInfo                         `json:"previous_period"`
	Metrics        Metrics                            `json:"metrics"`
	List           []PerformaStudioDetailItemResponse `json:"list"`
}

type PerformaStudioDetailItemResponse struct {
	AccountID   uint    `json:"account_id"`
	AccountName string  `json:"account_name"`
	GMV         int64   `json:"gmv"`
	Commission  int64   `json:"commission"`
	Ads         int64   `json:"ads"`
	Acos        float64 `json:"acos"`
	Roas        float64 `json:"roas"`
	Income      int64   `json:"income"`
}

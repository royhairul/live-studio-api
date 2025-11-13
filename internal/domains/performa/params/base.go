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

type PerformaMetricItem struct {
	GMV        int64   `json:"gmv"`
	Commission int64   `json:"commission"`
	Ads        int64   `json:"ads"`
	Income     int64   `json:"income"`
	Acos       float64 `json:"acos,omitempty"`
	Roas       float64 `json:"roas,omitempty"`
}
